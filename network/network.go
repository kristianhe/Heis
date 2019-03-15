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
	if err != nil 	{ fmt.Println("Error: ", localAddr, err) }
	socket, err := net.ListenUDP("udp", localAddr)
	if err != nil	{ fmt.Println("Error: ", localAddr, err) }
	return socket
}

func GetID(sender *net.UDPAddr) ID	{ return ID(sender.IP.String()) }

func GetIP() ID {
	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil	{ return "" }
	for _, interfaceAddrs := range interfaceAddrs {
		networkIP, ok := interfaceAddrs.(*net.IPNet)
		if ok && !networkIP.IP.IsLoopback() && networkIP.IP.To4() != nil	{ return networkIP.IP.String() }
	}
	return ""
}

func setDeadline(socket *net.UDPConn, t time.Time) {
	err := socket.SetReadDeadline(t.Add(time.Millisecond * 2000)) // Creates a deadline for just this socket, not all
	if err != nil && !err.(net.Error).Timeout() 	{ fmt.Println("Error: ", err) }
}

func listen(socket *net.UDPConn, incomming_information chan SimpleMessage, abort chan bool) {
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
								incomming_information <- SimpleMessage{GetID(sender), data[:receivedData]}
							} else if err != nil && !err.(net.Error).Timeout() {
								fmt.Println("Error: ", err)
							}
							time.Sleep(time.Millisecond * 10)
		}
	}
}

// TODO Siste argument er en slags sjekk. Finn ut hva denne gjør og lag en bra navn
func broadcast(socket *net.UDPConn, destination int, outgoing_information chan SimpleMessage, abort chan bool, isLocal bool) {
	address := GetIP()
	if !isLocal 	{ address = "255.255.255.255" }
	bcast_addr := fmt.Sprintf("%s:%d", address, destination)
	remote_addr := net.ResolveUDPAddr("udp", bcast_addr)
	if err != nil	{ fmt.Println("Error: ", err) }
	for {
		SimpleMessage := <-outgoing_information
		_, err := socket.WriteToUDP(SimpleMessage.Data, remote_addr)
		if err != nil {
			if isLocal {
				// Show error
				fmt.Println("Error: ", err)
				// Wait
				time.Sleep(time.Second * 1)
				// Reconnect
				broadcast(socket, destination, outgoing_information, abort, isLocal)
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// [Description here]
func Warden(read_from_slave chan SimpleMessage, write_to_slave chan SimpleMessage, abort chan bool) {
	socket := createSocket(masterPort)
	go listen(socket, read_from_slave, abort)
	broadcast(socket, slavePort, write_to_slave, abort, false)
	socket.Close()
}

// [Description here]
func Coordinator(read_from_master chan SimpleMessage, write_to_master chan SimpleMessage, abort chan bool) {
	socket := createSocket(slavePort)
	go listen(socket, read_from_master, abort)
	broadcast(socket, masterPort, write_to_master, abort, false)
	socket.Close()
}

// [Description here]
func BackupWarden(read_from_slave chan SimpleMessage, write_to_slave chan SimpleMessage, abort chan bool) {
	socket := createSocket(backupMasterPort)
	broadcast(socket, backupSlavePort, write_to_slave, abort, true)
	socket.Close()
}

// Continously listens to check if the master is alive
func BackupCoordinator(read_from_master chan SimpleMessage, write_to_master chan SimpleMessage, abort chan bool) {
	socket := createSocket(backupSlavePort)
	listen(socket, read_from_master, abort)
	socket.Close()
}
