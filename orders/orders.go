package orders

import (
	constants ".././common/constants"
	formats ".././common/formats"
	"../control"
	"../network"
	"../stateMachine"

	"fmt"
	"math"
	"sync"
	"time"
)

// Variables
var filename string = "[Orders] \t"
var mutex sync.Mutex
var isInserting bool
var orders []formats.Order
var ordersOffline []formats.Order

// [Description here]
func Handle(channel_poll_floor chan formats.Floor, channel_poll_order chan formats.Order, channel_write chan formats.SimpleMessage) {
	// Declare network messages
	var detailedMessageToSend formats.DetailedMessage
	var simpleMessageToSend formats.SimpleMessage
	for {
		select {
		// Called when we have a new order
		case currentOrder := <-channel_poll_order:
			// Prevent accessing the same memory at the same time
			if !isInserting {
				isInserting = true
				PrioritizeOrder(&currentOrder)
				// We do not want a duplicate order in the system
				if !CheckIfOrderExists(currentOrder) {
					// Define and send network messages
					detailedMessageToSend.Category = constants.MESSAGE_ORDER
					detailedMessageToSend.Order = currentOrder
					simpleMessageToSend.Data = network.EncodeMessage(detailedMessageToSend)
					channel_write <-simpleMessageToSend
					// Insert order into local array
					InsertOrder(currentOrder)
				}
				isInserting = false
			}
		// Called when we reach a new floor
		case currentFloor := <-channel_poll_floor:
			// Update floor
			stateMachine.SetFloor(currentFloor.Current)
			CheckIfOrdersAreCompleted(channel_write)
			RunActiveOrders()
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// [Description here]
func PrioritizeState(elevatorState int) int {
	if elevatorState == constants.STATE_IDLE || elevatorState == constants.STATE_DOOR_OPEN {
		return 1
	} else {
		return 0
	}
}

func PrioritizeFloor(currentFloor int, orderFloor int) int {
	return constants.FLOORS - int(math.Abs(float64(currentFloor - orderFloor)))
}

func PrioritizeDirection(elevatorState int, elevatorDirection int, orderDirection int) int {
	if elevatorState == constants.STATE_RUNNING {
		if elevatorDirection == orderDirection {
			return 1
		}
	}
	return 0
}

func PrioritizeOrders() {
	for index := range orders {
		PrioritizeOrder(&orders[index])
	}
	for index := range ordersOffline {
		PrioritizeOrder(&ordersOffline[index])
	}
	PrintOrders()
	fmt.Println(filename, "All orders are reprioritized.")
}

func PrioritizeOrder(order *formats.Order) {
	copy := GetOrder(*order)
	if copy.Category == constants.BUTTON_INSIDE {
		if copy.Elevator == stateMachine.GetConnectedIp() || copy.Elevator == constants.DEFAULT_IP {
			// Check network connection
			if stateMachine.IsConnected() {
				mutex.Lock()
				defer mutex.Unlock()
				order.Elevator = stateMachine.GetConnectedIp()
			}
			if !stateMachine.IsConnected() {
				mutex.Lock()
				defer mutex.Unlock()
				order.Elevator = constants.DEFAULT_IP
			}
		}
	}
	// Only prioritize orders pushed outside of elevators or local orders pushed inside
	if copy.Category == constants.BUTTON_OUTSIDE {
		var priority formats.Priority
		var priorities []formats.Priority
		var elevators []formats.Status = stateMachine.GetExternalElevators()
		// Find ID of local elevator
		priority.Elevator = network.GetIP()
		// Prioritize
		priority.Queue += PrioritizeState(stateMachine.GetState())
		priority.Queue += PrioritizeFloor(stateMachine.GetFloor(), copy.Floor)
		priority.Queue += PrioritizeDirection(stateMachine.GetState(), stateMachine.GetDirection(), copy.Direction)
		if stateMachine.GetState() == constants.STATE_EMERGENCY {
			priority.Queue -= 20
		}
		// Add priority
		priorities = append(priorities, priority)
		// External elevators
		if len(elevators) > 0 {
			for index := range elevators {
				// Fetch elevator
				priority.Elevator = elevators[index].Elevator
				priority.Queue = 0
				// Calculate priority
				priority.Queue += PrioritizeState(elevators[index].State)
				priority.Queue += PrioritizeFloor(elevators[index].Floor, copy.Floor)
				priority.Queue += PrioritizeDirection(elevators[index].State, elevators[index].Direction, copy.Direction)
				if elevators[index].State == constants.STATE_EMERGENCY {
					priority.Queue -= 20
				}
				// Add priority
				priorities = append(priorities, priority)
			}
		}
		// All priorities
		if len(priorities) > 0 {
			for index := range priorities {
				// If we have a bigger score
				if priorities[index].Queue > priority.Queue {
					// Change elevator
					priority.Elevator = priorities[index].Elevator
					priority.Queue = priorities[index].Queue
					// If we have the same score, compare IP addresses. Biggest IP get's the order.
				} else if priorities[index].Queue == priority.Queue {
					//Compare
					if priorities[index].Elevator > priority.Elevator {
						//Change elevator
						priority.Elevator = priorities[index].Elevator
						priority.Queue = priorities[index].Queue
					}
				}
			}
		}
		// Set target elevator to the most appropriate elevator
		mutex.Lock()
		order.Elevator = priority.Elevator
		order.Time = time.Now()
		mutex.Unlock()
	}
}

func InsertOrder(order formats.Order) {
	// Check if it is an order pushed inside or the correct direction outside 									// TODO hva søren betyr dette?
	if (order.Elevator == network.GetIP() && order.Category == constants.BUTTON_INSIDE) ||
		(order.Category == constants.BUTTON_OUTSIDE) {
		control.SetButtonLamp(order.Button, order.Floor, constants.ON)
	}
	// Add order
	mutex.Lock()
	orders = append(orders, order)
	mutex.Unlock()
	PrintOrder(order)
}

func InsertOfflineOrder(order formats.Order) {
	localOffline := GetOfflineOrders()
	isFound := false
	if len(localOffline) > 0 {
		for index := range localOffline {
			// Check if all fields are alike
			if IsOrdersEqual(localOffline[index], order) {
				isFound = true
			}
		}
	}
	// Add order
	if !isFound && order.Category == constants.BUTTON_INSIDE {
		mutex.Lock()
		ordersOffline = append(ordersOffline, order)
		mutex.Unlock()
		PrintOrder(order)
	}
}

func ResetOrderTimer(order formats.Order) {
	localOrders := GetOrders()
	for index := range localOrders {
		// Check if all fields are alike
		if IsOrdersEqual(localOrders[index], order) {
			orders[index].Time = time.Now()
		}
	}
}

func CountRelevantOrders() int {
	counter := 0
	localOrders := GetOrders()
	if len(localOrders) > 0 {
		for index := range localOrders {
			// Check if the order has something to do with local machine
			if localOrders[index].Elevator == network.GetIP() {
				counter++
			}
		}
	}
	return counter
}

func RunActiveOrders() {
	// Check if we have intiated an order
	if stateMachine.GetState() != constants.STATE_DOOR_OPEN {
		if CountRelevantOrders() > 0 {
			// Instantiate variables
			ordersOver := false
			ordersUnder := false
			localOrders := GetOrders()
			// Check if we have orders above or under current floor
			for index := range localOrders {
				// Check if it is a order pushed inside or the correct direction outside 				// TODO hva søren betyr denne kommentaren??
				if localOrders[index].Elevator == network.GetIP() {
					if localOrders[index].Floor > stateMachine.GetFloor() {
						ordersOver = true
					}
					if localOrders[index].Floor < stateMachine.GetFloor() {
						ordersUnder = true
					}
				}
			}
			// Run elevator in the correct direction
			if ordersOver && !(ordersUnder && stateMachine.GetDirection() == constants.DOWN) {
				control.SetMotorDir(constants.MOTOR_UP)
			} else if ordersUnder && !(ordersOver && stateMachine.GetDirection() == constants.UP) {
				control.SetMotorDir(constants.MOTOR_DOWN)
			}
			// If we don't have any orders over or under, change direction
			if !ordersOver && !ordersUnder {
				control.SwitchDir()
			}
		} else {
			// We have no new relevant orders
			control.SetMotorDir(constants.MOTOR_STOP)
			stateMachine.SetState(constants.STATE_IDLE)
		}
	}
}

func CheckIfOrderExists(order formats.Order) bool {
	localOrders := GetOrders()
	if len(localOrders) > 0 {
		for index := range localOrders {
			// Check if all fields are alike
			if IsOrdersEqual(localOrders[index], order) {
				return true
			}
		}
	}
	return false
}

func CheckIfOrdersAreCompleted(channel_write chan formats.SimpleMessage) {
	// Check if we have initiated an order
	if stateMachine.GetState() != constants.STATE_DOOR_OPEN {
		if CountRelevantOrders() > 0 {
			localOrders := GetOrders()
			for index := range localOrders {
				// Check if we are on the same floor as the order
				if localOrders[index].Floor == stateMachine.GetFloor() {
					// Check if it is a order pushed inside or the correct direction outside										// TODO igjen, hva betyr dette???
					if ((localOrders[index].Elevator == network.GetIP()) && (localOrders[index].Category == constants.BUTTON_INSIDE)) || (localOrders[index].Category == constants.BUTTON_OUTSIDE && localOrders[index].Direction == stateMachine.GetDirection()) {
						InitCompleteOrder(channel_write, localOrders[index])
						// Prevent panic
						break
					}
				}
			}
		}
	}
}

func InitCompleteOrder(channel_write chan formats.SimpleMessage, order formats.Order) {
	localOrders := GetOrders()
	for index := range localOrders {
		// Check if all fields are alike
		if IsOrdersEqual(localOrders[index], order) {
			// Open door for the specific elevator
			if order.Elevator == network.GetIP() {
				// Order-being processed sequence						// TODO hva betyr denne kommentaren????
				control.SetMotorDir(constants.MOTOR_STOP)
				stateMachine.SetState(constants.STATE_DOOR_OPEN)
				control.SetDoorLamp(constants.ON)
				timer := time.NewTimer(time.Second)
				// Finish order when timer is done
				go CompleteOrder(channel_write, order, timer)
			} else {
				control.SetButtonLamp(order.Button, order.Floor, constants.OFF)
				fmt.Println(filename, "Order completed.")
			}
			// Prevent reconnecting elevator double order handling			// TODO kan denne kommentaren forbedres??
			if !stateMachine.IsConnected() {
				InsertOfflineOrder(order)
			}
			// Remove order
			mutex.Lock()
			orders = append(orders[:index], orders[index+1:]...)
			mutex.Unlock()
			// Prevent panic
			break
		}
	}
}

func CompleteOrder(channel_write chan formats.SimpleMessage, order formats.Order, timer *time.Timer) {
	// When timer is finished
	<-timer.C // TODO hva er dette?
	// Order-complete sequence
	control.SetButtonLamp(order.Button, order.Floor, constants.OFF)
	control.SetDoorLamp(constants.OFF)
	stateMachine.SetState(constants.STATE_DOOR_CLOSED)
	// Network messages
	var detailedMessageToSend formats.DetailedMessage
	var simpleMessageToSend formats.SimpleMessage
	// Create and send message
	detailedMessageToSend.Category = constants.MESSAGE_FULFILLED
	detailedMessageToSend.Order = order
	simpleMessageToSend.Data = network.EncodeMessage(detailedMessageToSend)
	channel_write <- simpleMessageToSend
	fmt.Println(filename, "Order completed.")
}

// Requesting orders from other computers on the local network
func RequestOrders(channel_write chan formats.SimpleMessage) {
	fmt.Println(filename, "Requesting orders.")
	// Network messages
	var detailedMessageToSend formats.DetailedMessage
	var simpleMessageToSend formats.SimpleMessage
	// Create and send message
	detailedMessageToSend.Category = constants.MESSAGE_REQUEST
	simpleMessageToSend.Data = network.EncodeMessage(detailedMessageToSend)
	channel_write <- simpleMessageToSend
}

// Getters and setters
func GetOrder(order formats.Order) formats.Order {
	mutex.Lock()
	defer mutex.Unlock()
	// Make a full data copy
	copy := formats.Order{}
	copy = order
	return copy
}

func GetOrders() []formats.Order {
	mutex.Lock()
	defer mutex.Unlock()
	// Create a copy to prevent data race
	copy := make([]formats.Order, len(orders), len(orders))
	// Need to manually copy all variables (Library function "copy" will not do)
	for id, elem := range orders {
		copy[id] = elem
	}
	return copy
}

func GetOfflineOrders() []formats.Order {
	mutex.Lock()
	defer mutex.Unlock()
	// Create a copy to prevent data race
	copy := make([]formats.Order, len(ordersOffline), len(ordersOffline))
	// Need to manually copy all variables (Library function "copy" will not do)
	for id, elem := range ordersOffline {
		copy[id] = elem
	}
	return copy
}

func RemoveOfflineHistory() {
	mutex.Lock()
	defer mutex.Unlock()
	ordersOffline = ordersOffline[:0]
}

func SetOrders(list []formats.Order, channel_write chan formats.SimpleMessage) {
	for index := range list {
		isFound := false
		for offlineIndex := range ordersOffline {
			// Check if all fields are alike
			if IsOrdersEqual(list[index], ordersOffline[offlineIndex]) {
				isFound = true
			}
		}
		if !isFound {
			InsertOrder(list[index])
		}
	}
}

// Checking if orders are equal without checking timestamp
func IsOrdersEqual(order formats.Order, compareOrder formats.Order) bool { // TODO bør kanskje hete AreOrdersEqual?? Men lurt å ha isX-formatet
	if order.Category == constants.BUTTON_INSIDE {
		if order.Elevator != compareOrder.Elevator {
			return false
		}
	}
	if order.Category == compareOrder.Category &&
		order.Direction == compareOrder.Direction &&
		order.Floor == compareOrder.Floor &&
		order.Button == compareOrder.Button {
		return true
	} else {
		return false
	}
}

func PrintOrder(order formats.Order) {
	if order.Category == constants.BUTTON_INSIDE {
		fmt.Println(filename, "Order inside, IP: ", order.Elevator, ", floor: ", order.Floor)
	}
	if order.Category == constants.BUTTON_OUTSIDE {
		fmt.Print(filename, " Order outside, IP: ", order.Elevator, ", floor: ", order.Floor, ", direction: ")
		if order.Direction == constants.UP {
			fmt.Print(" up")
		} else {
			fmt.Print(" down")
		}
		fmt.Println()
	}
}

func PrintOrders() {
	localOrders := GetOrders()
	for index := range localOrders {
		PrintOrder(localOrders[index])
	}
}

func PrintPriority(local_priority formats.Priority) {
	fmt.Println(filename, " Elevator: ", local_priority.Elevator, ", count: ", local_priority.Queue)
}

func PrintPriorities(localPriorities []formats.Priority) {
	for index := range localPriorities {
		PrintPriority(localPriorities[index])
	}
}
