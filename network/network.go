package network

import (
	"../common"
	"fmt"
	"net"
	"time"
)

// Default ports
var masterPort int = 30012
var slavePort int = 30013
var backupMasterPort int = 30014
var backupSlavePort int = 30015

// Functions
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

func listen(socket *net.UDPConn, incomming_information chan formats.SimpleMessage, abort chan bool) {
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
								incomming_information <- formats.SimpleMessage{GetID(sender), data[:receivedData]}
							} else if err != nil && !err.(net.Error).Timeout() {
								fmt.Println("Error: ", err)
							}
							time.Sleep(time.Millisecond * 10)
		}
	}
}

// TODO Siste argument er en slags sjekk. Finn ut hva denne gjør og lag en bra navn
func broadcast(socket *net.UDPConn, destination int, outgoing_information chan formats.SimpleMessage, abort chan bool, isLocal bool) {
	address := GetIP()
	if !isLocal 	{ address = "255.255.255.255" }
	bcast_addr := fmt.Sprintf("%s:%d", address, destination)
	remote_addr := net.ResolveUDPAddr("udp", bcast_addr)
	if err != nil	{ fmt.Println("Error: ", err) }
	for {
		formats.SimpleMessage := <-outgoing_information
		_, err := socket.WriteToUDP(formats.SimpleMessage.Data, remote_addr)
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
func Warden(read_from_slave chan formats.SimpleMessage, write_to_slave chan formats.SimpleMessage, abort chan bool) {
	socket := createSocket(masterPort)
	go listen(socket, read_from_slave, abort)
	broadcast(socket, slavePort, write_to_slave, abort, false)
	socket.Close()
}

// [Description here]
func Coordinator(read_from_master chan formats.SimpleMessage, write_to_master chan formats.SimpleMessage, abort chan bool) {
	socket := createSocket(slavePort)
	go listen(socket, read_from_master, abort)
	broadcast(socket, masterPort, write_to_master, abort, false)
	socket.Close()
}

// [Description here]
func BackupWarden(read_from_slave chan formats.SimpleMessage, write_to_slave chan formats.SimpleMessage, abort chan bool) {
	socket := createSocket(backupMasterPort)
	broadcast(socket, backupSlavePort, write_to_slave, abort, true)
	socket.Close()
}

// Continously listens to check if the master is alive
// TODO kan fjerne argument nr 2? Det brukes ikke
func BackupCoordinator(read_from_master chan formats.SimpleMessage, write_to_master chan formats.SimpleMessage, abort chan bool) {
	socket := createSocket(backupSlavePort)
	listen(socket, read_from_master, abort)
	socket.Close()
}
