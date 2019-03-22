package cases

import (
	constants ".././common/constants"
	formats ".././common/formats"
	"../stateMachine"
	"../network"
	"../control"
	"../orders"

	"fmt"
	"time"
	"os"
	"os/signal"
)

var heartbeat = time.Now()

// TODO Vi har de to neste funksjonene i elevio også, diskuter hvor de bør være

func PollFloor(floorChannel chan formats.Floor) {
	for {
		polledFloor := control.GetFloorSignal()
		if polledFloor != constants.INVALID {
			var newFloor formats.Floor
			newFloor.Current = polledFloor
			control.SetFloorIndicator(polledFloor)
			floorChannel <-newFloor
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func PollOrder(orderChannel chan formats.Order) {
	for {
		// Get registered button and floor requested
		polledOrderButton, polledOrderFloor := CheckRequestedOrders()
		// Only go through if not invalid button 																			// TODO sjekk kommentarene her....
		if polledOrderButton != constants.INVALID && polledOrderFloor != constants.INVALID {
			var newOrder formats.Order
			newOrder.Elevator = network.GetIP()
			if polledOrderButton == constants.BUTTON_INSIDE {
				newOrder.Category 	= constants.BUTTON_INSIDE
			}
			if polledOrderButton == constants.BUTTON_UP {
				newOrder.Category 	= constants.BUTTON_OUTSIDE
				newOrder.Direction 	= constants.UP
			}
			if polledOrderButton == constants.BUTTON_DOWN {
				newOrder.Category 	= constants.BUTTON_OUTSIDE
				newOrder.Direction 	= constants.DOWN
			}
			// Register the new order
			newOrder.Floor = polledOrderFloor
			newOrder.Button = polledOrderButton
			newOrder.Time = time.Now()
			orderChannel <-newOrder
		}
		time.Sleep(time.Millisecond * 10)
	}
}


// Sends a message to the local master every half second										// TODO hva betyr egentlig 'local master'?
func Heartbeater(backupChannel_write chan formats.SimpleMessage) {
	// Declare network messages
	var detailedMessageToSend formats.DetailedMessage
	var simpleMessageToSend formats.SimpleMessage
	for {
		if stateMachine.IsMaster() {
			// Define and send network messages
			detailedMessageToSend.Category = constants.MESSAGE_HEARTBEAT
			simpleMessageToSend.Data = network.EncodeMessage(detailedMessageToSend)
			backupChannel_write <-simpleMessageToSend
		}
		time.Sleep(time.Millisecond * 500)
	}
}

// Checks if we have received a heartbeat from the local master. If it takes longer than three seconds, spawn as new master.				// TODO hva var denne lokale masteren igjen?
func CheckHeartbeat(channel_abort chan bool, channel_init_master chan bool) {
	fmt.Println("CheckHeartbeat has started")
	for {
		if !stateMachine.IsMaster() {
			// Calculate time
			elapsedTime := time.Since(heartbeat)
			elapsedTime = (elapsedTime + time.Second/2) / time.Second															// TODO forstå denne....
			if elapsedTime > 3 {
				// Spawn as master
				channel_init_master <-true																	// TODO har endret litt på rekkefølgen her... Håper ikke det er utslagsgivende
				// End this goroutine
				channel_abort <-true
				fmt.Println("Did not receive heartbeat. Closing socket and rebooting as master.")
				return
			}
		} else {
			// No longer master, so break the loop
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// Checks if we have received a heartbeat form the local backup. Update timestamp in receive.							// TODO Local backup? Har vi ekstern backup også??
func CheckBackupHeartbeat(channel_init_master chan bool, backupChannel_read chan formats.SimpleMessage) {					// TODO Sjekk om denne tolkningen av koken er rett...
	fmt.Println("CheckBackupHeartbeat has started")
	// Declare network messages
	var detailedMessageReceived formats.DetailedMessage
	var simpleMessageReceived formats.SimpleMessage
	for {
		if !stateMachine.IsMaster() {
			select {
			case simpleMessageReceived = <-backupChannel_read:
				// Get message and decode
				detailedMessageReceived = network.DecodeMessage(simpleMessageReceived.Data)
				// Check heartbeat from main process on this computer 														// TODO gjør noe med denne kommentaren.. Hva betyr main process
				if detailedMessageReceived.Category == constants.MESSAGE_HEARTBEAT && simpleMessageReceived.Address == network.GetIP()  { resetHeartbeat() }
			}
		} else {
			// No longer master, so break the loop
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func resetHeartbeat()  { heartbeat = time.Now() }

// [Description here]
func BroadcastToNetwork(channel_write chan formats.SimpleMessage) {
	// Declare network messages
	var detailedMessageToSend formats.DetailedMessage
	var simpleMessageToSend formats.SimpleMessage
	for {
		if stateMachine.IsMaster() {
			// Create status messages
			detailedMessageToSend.Category 			= constants.MESSAGE_STATUS
			detailedMessageToSend.Status.Elevator 	= network.GetIP()
			detailedMessageToSend.Status.State 		= stateMachine.GetState()
			detailedMessageToSend.Status.Floor 		= stateMachine.GetFloor()
			detailedMessageToSend.Status.Direction 	= stateMachine.GetDirection()
			detailedMessageToSend.Status.Time 		= time.Now()
			// Encode the message and send
			simpleMessageToSend.Data = network.EncodeMessage(detailedMessageToSend)
			channel_write <-simpleMessageToSend
		}
		time.Sleep(time.Millisecond * 500)
	}
}

// Listens to the network, receives messagees from the other elevators and does actions according to the message category		// TODO bør denne heller være i network-modulen
func ListenToNetwork(channel_read chan formats.SimpleMessage, channel_write chan formats.SimpleMessage) {
	// Declare network messages
	var detailedMessageReceived formats.DetailedMessage
	var simpleMessageReceived formats.SimpleMessage
	for {
		if stateMachine.IsMaster() {
			select {
			case simpleMessageReceived = <-channel_read:
				// Get message and decode
				detailedMessageReceived = network.DecodeMessage(simpleMessageReceived.Data)
				// Ignore messages sent from this computer
				if simpleMessageReceived.Address != network.GetIP() {
					// Case structure for all message categories
					switch detailedMessageReceived.Category {
					case constants.MESSAGE_STATUS:
						stateMachine.UpdateExternalElevator(detailedMessageReceived.Status)
						break
					case constants.MESSAGE_ORDER:
						// Order received from another elevator
						fmt.Println("Order received.")
						// Check if the message wasn't sent from this computer
						if !orders.CheckIfOrderExists(detailedMessageReceived.Order) { orders.InsertOrder(detailedMessageReceived.Order) }
						break
					case constants.MESSAGE_FULFILLED:
						// Order fulfilled by another elevator
						fmt.Println("Order already fulfilled.")
						orders.InitCompleteOrder(channel_write, detailedMessageReceived.Order)													// TODO hva gjør egentlig denne??
						break
					case constants.MESSAGE_ORDERS:
						fmt.Println("Order list received.")
						// Check if this is the intended destination
						if detailedMessageReceived.Orders.Elevator == network.GetIP()  { orders.SetOrders(detailedMessageReceived.Orders.List, channel_write) }
						// Prevent reconnectiong elevator double order handling 											// TODO hva i all verden betyr dette???
						orders.RemoveOfflineHistory()
						break
					case constants.MESSAGE_REQUEST:
						// Request order list
						fmt.Println("Elevator", simpleMessageReceived.Address, "is requesting orders.")
						// Declare network messages
						var detailedMessageToSend formats.DetailedMessage
						var simpleMessageToSend formats.SimpleMessage
						// Define network messages and send them
						detailedMessageToSend.Category 			= constants.MESSAGE_ORDERS
						detailedMessageToSend.Orders.Elevator 	= simpleMessageReceived.Address
						detailedMessageToSend.Orders.List 		= orders.GetOrders()
						simpleMessageToSend.Data 				= network.EncodeMessage(detailedMessageToSend)
						channel_write <-simpleMessageReceived
						fmt.Println("Orders sent.")
						break
					case constants.MESSAGE_REPRIORITIZE:
						// Order fulfilled by another elevator
						fmt.Println("Reprioritizing.")
						orders.PrioritizeOrders()
						break
					}
				}
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// [Description here]
func CheckStatus(channel_write chan formats.SimpleMessage) {											// TODO kan sannsynligvis finne et bedre navn på denne funksjonen
	// Declare network messages
	var detailedMessageToSend formats.DetailedMessage
	var simpleMessageToSend formats.SimpleMessage
	for {
		// If a disconnected computer gets connection, it can request orders from external elevators.
		if !stateMachine.IsConnected() && network.GetIP() != constants.DEFAULT_IP {
			stateMachine.SetConnection(constants.CONNECTED)
			stateMachine.SetConnectedIp(network.GetIP())
			time.Sleep(time.Second)
			orders.PrioritizeOrders()
			orders.RequestOrders(channel_write)
		}
		// If a connected computer looses connection, it only works locally.
		if stateMachine.IsConnected() && network.GetIP() == constants.DEFAULT_IP {
			stateMachine.SetConnection(constants.DISCONNECTED)
			time.Sleep(time.Second)
			orders.PrioritizeOrders()
		}
		// [Descripton here]
		localElevators := stateMachine.GetExternalElevators()
		for elevatorIndex := range localElevators {
			// Fetch spesific elevator (for simpler syntax)
			elevator := localElevators[elevatorIndex]
			// Calculate time
			elapsedTime := time.Since(elevator.Time)
			elapsedTime = (elapsedTime + time.Second/2) / time.Second
			// Check for timeout
			if elapsedTime > 3 {
				stateMachine.RemoveExternalElevator(elevator)
				orders.PrioritizeOrders()
			}
		}
		// [Descripton here]
		if stateMachine.GetState() != constants.STATE_EMERGENCY {
			localOrders := orders.GetOrders()
			for ordersIndex := range localOrders {
				// Fetch spesific order (for simpler syntax)
				order := localOrders[ordersIndex]
				// Calculate time
				elapsedTime := time.Since(order.Time)
				elapsedTime = (elapsedTime + time.Second/2) / time.Second
				// Check for timeout
				if order.Elevator == network.GetIP() && elapsedTime > 20 {
					stateMachine.SetState(constants.STATE_EMERGENCY)
					orders.PrioritizeOrders()
					// Give time to let external computers to receive the status message
					time.Sleep(time.Second)
					// Define the network message and send it
					detailedMessageToSend.Category = constants.MESSAGE_REPRIORITIZE
					simpleMessageToSend.Data = network.EncodeMessage(detailedMessageToSend)
					channel_write <-simpleMessageToSend
					fmt.Println("Request for repriorization sent.")
				}
			}
		}
		time.Sleep(time.Second)
	}
}

// [Description here]
func CheckRequestedOrders() (int, int) {
	for floor := 0; floor < constants.FLOORS; floor++ {
		for button := 0; button < constants.BUTTONS; button++ {
			if control.GetButtonSignal(button, floor) == true { // TODO stor sjanse for at det er noe feil med denne. Den er omgjort fra en int til en bool
				return button, floor
			}
		}
	}
	return constants.INVALID, constants.INVALID
}

// Handles exits by Ctrl+C termination
func ExitHandler() {
	// This channel is called when the program is terminated
	signalChannel := make(chan os.Signal, 10)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	// Stop the elevator
	control.Stop()
	fmt.Println("The program is killed! Elevator has stopped.")
	// Do last actions and wait for all write operations to end
	os.Exit(0)
}
