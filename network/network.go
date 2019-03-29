package network

import (
	formats "../common/formats"

	"encoding/json"
	"fmt"
	"net"
	"time"
)

var filename string = "[Network] \t"

// Default ports
var slavePort int = 20000
var masterPort int = 20001
var backupSlavePort int = 20002
var backupMasterPort int = 20003

func createSocket(port int) *net.UDPConn {
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error:", localAddr, err)
	}
	socket, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Println("Error:", localAddr, err)
	}
	return socket
}

func setDeadline(socket *net.UDPConn, t time.Time) {
	// Creates a deadline for just this socket, not all
	err := socket.SetReadDeadline(t.Add(time.Millisecond * 2000))
	if err != nil && !err.(net.Error).Timeout() {
		fmt.Println("Error: ", err)
	}
}

func GetIP() formats.ID {
	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, interfaceAddrs := range interfaceAddrs {
		networkIP, ok := interfaceAddrs.(*net.IPNet)
		if ok && !networkIP.IP.IsLoopback() && networkIP.IP.To4() != nil {
			return formats.ID(networkIP.IP.String())
		}
	}
	return ""
}

func GetID(sender *net.UDPAddr) formats.ID {
	 return formats.ID(sender.IP.String())
 }

 func EncodeMessage(msg formats.DetailedMessage) []byte {
 	result, err := json.Marshal(msg)
 	if err != nil {
 		fmt.Println("Error: ", err)
 	}
 	return result
 }

 func DecodeMessage(b []byte) formats.DetailedMessage {
 	var result formats.DetailedMessage
 	err := json.Unmarshal(b, &result)
 	if err != nil {
 		fmt.Println("Error: ", err)
 	}
 	return result
 }

func listen(socket *net.UDPConn, incomming_information chan formats.SimpleMessage, abort chan bool) {
	for {
		select {
		case <-abort:
			socket.Close()
			return
		default:
			 // Deadline for the listener
			setDeadline(socket, time.Now())
			data := make([]byte, 2048)
			receivedData, sender, err := socket.ReadFromUDP(data)
			if err == nil {
				incomming_information <-formats.SimpleMessage{GetID(sender), data[:receivedData]}
			} else if err != nil && !err.(net.Error).Timeout() {
				fmt.Println("Error: ", err)
			}
			time.Sleep(time.Millisecond * 10)

		}
	}
}

func broadcast(socket *net.UDPConn, destination int, outgoing_information chan formats.SimpleMessage, abort chan bool, isLocal bool) {
	// If broadcast is set to local
	address := GetIP()
	// Else
	if !isLocal {
		address = "255.255.255.255"
	}
	bcast_addr := fmt.Sprintf("%s:%d", address, destination)
	remote_addr, err := net.ResolveUDPAddr("udp", bcast_addr)
	if err != nil {
		fmt.Println("Error: ", err)
	}
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

// Deals with communication from the master to the slave
func MasterCoordinator(read_from_slave chan formats.SimpleMessage, write_to_slave chan formats.SimpleMessage, abort chan bool) {
	socket := createSocket(masterPort)
	go listen(socket, read_from_slave, abort)
	broadcast(socket, slavePort, write_to_slave, abort, false)
	socket.Close()
}

//  Deals with communication from the slave to the master
func SlaveCoordinator(read_from_master chan formats.SimpleMessage, write_to_master chan formats.SimpleMessage, abort chan bool) {
	socket := createSocket(slavePort)
	go listen(socket, read_from_master, abort)
	broadcast(socket, masterPort, write_to_master, abort, false)
	socket.Close()
}

// Coordinates the channel where heartbeats are sent from the backup to the master
func BackupCoordinator(write_to_slave chan formats.SimpleMessage, abort chan bool) {
	socket := createSocket(backupMasterPort)
	broadcast(socket, backupSlavePort, write_to_slave, abort, true)
	socket.Close()
}

// Listens to the channel where heartbeats are sent from the backup to the master
func BackupListener(read_from_master chan formats.SimpleMessage, abort chan bool) {
	socket := createSocket(backupSlavePort)
	listen(socket, read_from_master, abort)
	socket.Close()
}
