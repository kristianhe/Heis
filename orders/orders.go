package orders

import (
	"../control"
	"../network"
	"../stateMachine"
	"../tools"
	"fmt"
	"math"
	"sync"
	"time"

)

//Variables
var filename string = "Orders -"
var mutex sync.Mutex
var inserting bool
var orders []tools.Order
var orders_offline []tools.Order
var check_counter = 0

func Handle(channel_poll_floor chan tools.Floor, channel_poll_order chan tools.Order, channel_write chan tools.Packet) {
	//Network messages
	var message_send tools.Message
	var packet_send tools.Packet
	for {
		select {
		//Called when we have a new order
		case current_order := <-channel_poll_order:
			//Prevent accessing the same memory at the same time
			if !inserting {
				inserting = true
				PrioritizeOrder(&current_order)
				//We we not want a duplicate order in the system
				if !CheckOrderExists(current_order) {
					//Create and send message
					message_send.Category = tools.MESSAGE_ORDER
					message_send.Order = current_order
					packet_send.Data = network.EncodeMessage(message_send)
					channel_write <- packet_send
					//Insert order into local array
					InsertOrder(current_order)
				}
				inserting = false
			}
		//Called when we reach a new floor
		case current_floor := <-channel_poll_floor:
			//Update floor
			states.SetFloor(current_floor.Current)
			CheckOrdersCompleted(channel_write)
			RunActiveOrders()
			//Check if we are at bottom
			if states.GetFloor() == tools.FLOOR_FIRST {
				driver.DirectionUp()
			}
			//Check if we are at top
			if states.GetFloor() == tools.FLOOR_LAST {
				driver.DirectionDown()
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func PrioritizeState(elevator_state int) int {
	if elevator_state == tools.STATE_IDLE || elevator_state == tools.STATE_DOOR_OPEN  {
		return 1
	}
	return 0
}

func PrioritizeFloor(current_floor int, order_floor int) int {
	return tools.FLOORS - int(math.Abs(float64(current_floor-order_floor)))
}

func PrioritizeDirection(elevator_state int, elevator_direction int, order_direction int) int {
	if elevator_state == tools.STATE_RUNNING {
		if elevator_direction == order_direction {
			return 1
		}
	}
	return 0
}

func PrioritizeOrders() {
	for index := range orders {
		PrioritizeOrder(&orders[index])
	}

	for index := range orders_offline {
		PrioritizeOrder(&orders_offline[index])
	}
	PrintOrders()
	fmt.Println(filename, "All orders reprioritized")
}

func PrioritizeOrder(order *tools.Order) {
	copy := GetOrder(*order)
	if copy.Category == tools.BUTTON_INSIDE {
		if copy.Elevator == states.GetConnectedIp() || copy.Elevator == tools.DEFAULT_IP {
			//Check network connection
			if states.IsConnected() {
				mutex.Lock()
				order.Elevator = states.GetConnectedIp()
				mutex.Unlock()
			}
			//If disconnected
			if !states.IsConnected() {
				mutex.Lock()
				order.Elevator = tools.DEFAULT_IP
				mutex.Unlock()
			}
		}
	}
	//Only prioritize orders pushed outside of elevators or local orders pushed inside
	if copy.Category == tools.BUTTON_OUTSIDE {
		var priority tools.Priority
		var priorities []tools.Priority
		var elevators []tools.Status = states.GetExternalElevators()
		//Local elevator
		priority.Elevator = network.GetMachineID()
		//Prioritize
		priority.Count += PrioritizeState(states.GetState())
		priority.Count += PrioritizeFloor(states.GetFloor(), copy.Floor)
		priority.Count += PrioritizeDirection(states.GetState(), states.GetDirection(), copy.Direction)
		if(states.GetState() == tools.STATE_EMERGENCY){
			priority.Count -= 20 
		}
		//Add priority
		priorities = append(priorities, priority)
		//External elevators
		if len(elevators) > 0 {
			for index := range elevators {
				//Fetch elevator
				priority.Elevator = elevators[index].Elevator
				priority.Count = 0
				//Calculate priority
				priority.Count += PrioritizeState(elevators[index].State)
				priority.Count += PrioritizeFloor(elevators[index].Floor, copy.Floor)
				priority.Count += PrioritizeDirection(elevators[index].State, elevators[index].Direction, copy.Direction)
				if(elevators[index].State == tools.STATE_EMERGENCY){
					priority.Count -= 20 
				}
				//Add priority
				priorities = append(priorities, priority)

			}
		}
		//All priorities
		if len(priorities) > 0 {
			for index := range priorities {
				//If we have a bigger score
				if priorities[index].Count > priority.Count {
					//Change elevator
					priority.Elevator = priorities[index].Elevator
					priority.Count = priorities[index].Count
				//If we have the same score - compare IP's. Biggest IP gets order.
				} else if priorities[index].Count == priority.Count {
					//Compare
					if priorities[index].Elevator > priority.Elevator {
						//Change elevator 
						priority.Elevator = priorities[index].Elevator
						priority.Count = priorities[index].Count
					}
				}

			}
		}
		//Set target elevator to the elevator most appropriate
		mutex.Lock()
		order.Elevator = priority.Elevator
		order.Time = time.Now()
		mutex.Unlock()
	}
}

func InsertOrder(order tools.Order) {
	//Check if it is a order pushed inside or the correct direction outside
	if (order.Elevator == network.GetMachineID() && order.Category == tools.BUTTON_INSIDE) || 
		(order.Category == tools.BUTTON_OUTSIDE) {driver.SetButtonLamp(order.Button, order.Floor, tools.ON)
	}
	//Add order
	mutex.Lock()
	orders = append(orders, order)
	mutex.Unlock()
	//Successfully added
	PrintOrder(order)
}

func InsertOfflineOrder(order tools.Order) {
	local_offline := GetOfflineOrders()
	found := false
	if len(local_offline) > 0 {
		for index := range local_offline {
			//Check if all fields are alike
			if IsOrdersEqual(local_offline[index], order) {
				found = true
			}
		}
	}
	//Add order
	if !found && order.Category == tools.BUTTON_INSIDE {
		mutex.Lock()
		orders_offline = append(orders_offline, order)
		mutex.Unlock()
		//Successfully added
		PrintOrder(order)
	}
}

func ResetOrderTimer(order tools.Order){
	local_orders := GetOrders()
	for index := range local_orders {
		//Check if all fields are alike
		if IsOrdersEqual(local_orders[index], order) {
			orders[index].Time = time.Now()
		}
	}
}

func CountRelevantOrders() int {
	counter := 0
	local_orders := GetOrders()
	if len(local_orders) > 0 {
		for index := range local_orders {
			//Check if the order has something to do with local machine
			if local_orders[index].Elevator == network.GetMachineID() {
				//Count
				counter++
			}
		}
	}
	return counter
}

func RunActiveOrders() {
	//Check if we have intiated an order
	if states.GetState() != tools.STATE_DOOR_OPEN {
		if CountRelevantOrders() > 0 {
			//Initiate variables
			orders_over := false
			orders_under := false
			local_orders := GetOrders()
			//Check if we have orders above or under current floor
			for index := range local_orders {
				//Check if it is a order pushed inside or the correct direction outside
				if local_orders[index].Elevator == network.GetMachineID() {
					if local_orders[index].Floor > states.GetFloor() {
						orders_over = true
					}
					if local_orders[index].Floor < states.GetFloor() {
						orders_under = true
					}
				}
			}
			//Run elevator in the correct direction
			if orders_over && !(orders_under && states.GetDirection() == tools.DOWN) {
				driver.RunUp()
			} else if orders_under && !(orders_over && states.GetDirection() == tools.UP) {
				driver.RunDown()
			}
			//If we dont have orders over or under, change direction
			if !orders_over && !orders_under {
				driver.DirectionSwitch()
			}
		} else {
			//We have no new relevant orders.
			driver.Stop()
			states.SetState(tools.STATE_IDLE)
		}
	}
}

func CheckOrderExists(order tools.Order) bool {
	local_orders := GetOrders()
	if len(local_orders) > 0 {
		for index := range local_orders {
			//Check if all fields are alike
			if IsOrdersEqual(local_orders[index], order) {
				return true
			}
		}
	}
	return false
}

func CheckOrdersCompleted(channel_write chan tools.Packet) {
	//Check if we have intiated an order
	if states.GetState() != tools.STATE_DOOR_OPEN {
		if CountRelevantOrders() > 0 {
			local_orders := GetOrders()
			for index := range local_orders {
				//Check if we are on the same floor as the order
				if local_orders[index].Floor == states.GetFloor() {
					//Check if it is a order pushed inside or the correct direction outside
					if (local_orders[index].Elevator == network.GetMachineID()) && (local_orders[index].Category == tools.BUTTON_INSIDE) || (local_orders[index].Category == tools.BUTTON_OUTSIDE 
						&& local_orders[index].Direction == states.GetDirection()) {InitCompleteOrder(channel_write, local_orders[index])
						//Prevent panic
						break
					}
				}
			}
		}
	}
}

func InitCompleteOrder(channel_write chan tools.Packet, order tools.Order) {
	local_orders := GetOrders()
	for index := range local_orders {
		//Check if all fields are alike
		if IsOrdersEqual(local_orders[index], order) {
			//Open door for the specific elevator
			if order.Elevator == network.GetMachineID() {
				//Order-being processed sequence
				driver.Stop()
				states.SetState(tools.STATE_DOOR_OPEN)
				driver.SetDoorLamp(tools.ON)
				timer := time.NewTimer(time.Second)
				//Finish order when timer is done
				go CompleteOrder(channel_write, order, timer)

			} else {
				driver.SetButtonLamp(order.Button, order.Floor, tools.OFF)
				fmt.Println(filename, "Order completed")
			}
			//Prevent reconnecting elevator double order handling
			if !states.IsConnected() {
				InsertOfflineOrder(order)
			}
			//Remove order
			mutex.Lock()
			orders = append(orders[:index], orders[index+1:]...)
			mutex.Unlock()
			//Prevent panic
			break
		}
	}
}

func CompleteOrder(channel_write chan tools.Packet, order tools.Order, timer *time.Timer) {
	//When timer is finished
	<-timer.C
	//Order-complete sequence
	driver.SetButtonLamp(order.Button, order.Floor, tools.OFF)
	driver.SetDoorLamp(tools.OFF)
	states.SetState(tools.STATE_DOOR_CLOSED)
	//Network messages
	var message_send tools.Message
	var packet_send tools.Packet
	//Create and send message
	message_send.Category = tools.MESSAGE_FULFILLED
	message_send.Order = order
	packet_send.Data = network.EncodeMessage(message_send)
	channel_write <- packet_send
	fmt.Println(filename, "Order completed")
}

//Requesting orders from other computers on the local network
func RequestOrders(channel_write chan tools.Packet) {
	fmt.Println(filename, "Requesting orders")
	//Network messages
	var message_send tools.Message
	var packet_send tools.Packet
	//Create and send message
	message_send.Category = tools.MESSAGE_REQUEST
	packet_send.Data = network.EncodeMessage(message_send)
	channel_write <- packet_send
}

//Getters and setters
func GetOrder(order tools.Order) tools.Order {
	mutex.Lock()
	defer mutex.Unlock()
	//Make a full data copy
	o := tools.Order{}
	o = order
	return o
}

func GetOrders() []tools.Order {
	mutex.Lock()
	defer mutex.Unlock()
	//Create a copy - preventing data race
	o := make([]tools.Order, len(orders), len(orders))
	//Need to manually copy all variables - Library "copy" function will not work
	for id, elem := range orders {
		o[id] = elem
	}
	return o
}

func GetOfflineOrders() []tools.Order {
	mutex.Lock()
	defer mutex.Unlock()
	//Create a copy - preventing data race
	o := make([]tools.Order, len(orders_offline), len(orders_offline))
	//Need to manually copy all variables - Library "copy" function will not work
	for id, elem := range orders_offline {
		o[id] = elem
	}
	return o
}

func RemoveOfflineHistory() {
	mutex.Lock()
	orders_offline = orders_offline[:0]
	mutex.Unlock()
}

func SetOrders(list []tools.Order, channel_write chan tools.Packet) {
	for index := range list {
		found := false
		for offline_index := range orders_offline {
			//Check if all fields are alike
			if IsOrdersEqual(list[index], orders_offline[offline_index]) {
				found = true
			}
		}
		if !found {
			InsertOrder(list[index])
		}
	}
}

//Checking if orders are equal without checking timestamp
func IsOrdersEqual(order tools.Order, compare_order tools.Order) bool{
	if(order.Category == tools.BUTTON_INSIDE){
		if(order.Elevator != compare_order.Elevator){
			return false
		}
	}
	if( order.Category  == compare_order.Category &&
		order.Direction == compare_order.Direction &&
		order.Floor     == compare_order.Floor &&
		order.Button    == compare_order.Button){
		return true
	}
	return false
}

//Printing functionality
func PrintOrders() {
	local_orders := GetOrders()
	for index := range local_orders {
		PrintOrder(local_orders[index])
	}
}

func PrintOrder(order tools.Order) {
	if order.Category == tools.BUTTON_INSIDE {
		fmt.Println(filename, "Order inside, ip: ", order.Elevator, ",floor: ", order.Floor)
	}
	if order.Category == tools.BUTTON_OUTSIDE {
		fmt.Print(filename, " Order outside, ip: ", order.Elevator, ",floor: ", order.Floor, ",direction: ")
		if order.Direction == tools.UP {
			fmt.Print(" up")
		} else {
			fmt.Print(" down")
		}
		fmt.Println()
	}
}

func PrintPriority(local_priority tools.Priority) {
	fmt.Println(filename, " Elevator: ", local_priority.Elevator, ", count: ", local_priority.Count)
}

func PrintPriorities(local_priorities []tools.Priority) {
	for index := range local_priorities {
		PrintPriority(local_priorities[index])
	}
}

