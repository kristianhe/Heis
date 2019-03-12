package network

import (
	"fmt"
	"net"
	"time"
)

// Default ports
var masterPort int = 30012
var slavePort int = 30013
var backupMasterPort int = 30014
var backupSlavePort int = 30015

// Standard message formats used in the network communication
type ID string

type SimpleMessage struct {
	Address ID
	Data    []byte
}

type DetailedMessage struct {
	Category  int // Message category ?? heller en enum?
	Heartbeat Heartbeat
	Status    Status
	Order     Order
	OrderList OrderList
}

type Status struct {
	Elevator  ID
	State     int
	Floor     int
	Direction int
	Priority  int // Bruke Priority-structen her?
}

type Order struct {
	Category  string // Message category ??
	Elevator  ID
	Direction int
	Floor     int
	Button    int
	time      time.Time
}

type OrderList struct {
	Elevator ID
	List     []Order
}

type Floor struct {
	Current int
	Status  int // Moving, idle etc. Bruke enum her i stedet for int?
}

type Priority struct {
	Elevator ID
	Queue    int // Place in queue
}

type Heartbeat struct {
	Count int
}

func createSocket(port int) *net.UDPConn {

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error: ", localAddr, err)
	}
	socket, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Println("Error: ", localAddr, err)
	}

	return socket
}

func getID(sender *net.UDPAddr) ID { return ID(sender.IP.String()) }

func setDeadline(socket *net.UDPConn, t time.Time) {

	err := socket.SetReadDeadline(t.Add(time.Millisecond * 2000)) // Creates a deadline for just this socket, not all
	if err != nil && !err.(net.Error).Timeout() {
		fmt.Println("Error: ", err)
	}

}

// TODO det er noe muffins med denne funksjonen.... får merkelige errormeldinger...
/*
func listenToSlaves(socket *net.UDPConn, received_information chan SimpleMessage, abort chan bool) {

	for {
		select {
			case <-abort:
							socket.Close()
							return
			default:
							setDeadline(socket, time.Now())	// Deadline for the listener
							data := make([]byte, 2048)
							receivedData, sender, err := socket.ReadFromUDP(data)
							if err == nil {

								received_information <- SimpleMessage{getID(sender), data[:receivedData]}

							} else if err != nil && !err.(net.Error).Timeout() {

								fmt.Println("Error: ", err)

							}
							time.Sleep(time.Millisecond * 10)

		}
	}

}
*/

// Siste argument er en slags sjekk, vet ikke om vi trenger den?
func broadcast(socket *net.UDPConn, destination int, information_to_send chan SimpleMessage, abort chan bool, some_error_check bool) {
	return
}

/*
// (MasterWorker)
func Warden(read_from_slave chan SimpleMessage, write_to_slave chan SimpleMessage, abort chan bool) {
	return
}

// (ClientWorker)
func Coordinator(read_from_master chan SimpleMessage, write_to_master chan SimpleMessage, abort chan bool) {
	return
}
*/

// (MasterBackupWorker)
func BackupWarden(read_from_slave chan SimpleMessage, write_to_slave chan SimpleMessage, abort chan bool) {

	socket := createSocket(backupMasterPort)
	broadcast(socket, backupSlavePort, write_to_slave, abort, true)
	socket.Close()

}

// (ClientBackupWorker) Continously listens to check if the master alive
func BackupCoordinator(read_from_master chan SimpleMessage, write_to_master chan SimpleMessage, abort chan bool) {

	socket := createSocket(backupSlavePort)
	//listenToSlaves(socket, read_from_master, abort)
	socket.Close()

}

/*

// TODO make a better name for this function
func NetworkFunc() {

	// Can set personal ID using two alternatives,
	// Alternative 1:
	// Choose by writing `go run main.go -id=your_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// Alternative 2:
	// Preset by local IP address + process ID
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// Channels for peer status
	peerUpdateCh := make(chan peers.PeerUpdate) // Peer status (who is currently active)
	peerTxEnable := make(chan bool)             // Enable/disable transmission

	// Start corresponding goroutines
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// Local channels for sending and receiving the message packet
	msgTx := make(chan Msg)
	msgRx := make(chan Msg)

	// Start broadcasting
	go bcast.Transmitter(16569, msgTx)
	go bcast.Receiver(16569, msgRx)

	// Send a message every second.
	go func() {
		helloMsg := Msg{"Hello from " + id, 0, 0}
		for {
			helloMsg.Iter++
			msgTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	// While true
	fmt.Println()
	fmt.Println("----------------------")
	fmt.Println("NETWORK UP AND RUNNING")
	fmt.Println("----------------------")
	fmt.Println()
	for {
		select {
		case peerUpdate := <-peerUpdateCh:
			fmt.Println("----------------------------")
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peerUpdate.Peers)
			fmt.Printf("  New:      %q\n", peerUpdate.New)
			fmt.Printf("  Lost:     %q\n", peerUpdate.Lost)
			fmt.Println("----------------------------")

		case newMessage := <-msgRx:
			fmt.Printf("Received: %#v\n", newMessage)

		}
	}
}
*/
