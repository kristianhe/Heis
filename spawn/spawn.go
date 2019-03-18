package spawn

import (
	"fmt"
	"os/exec"
	//"../common" // Denne bare fjernes når jeg lagrer filen...?
	"../control"
	"../network"
)

// Main channels
var channel_write = make(chan formats.SimpleMessage)
var channel_read = make(chan formats.SimpleMessage)

// Backup channels
var backupChannel_write = make(chan formats.SimpleMessage)
var backupChannel_read = make(chan formats.SimpleMessage)

// Toggle channels
var channel_init_master = make(chan bool)
var channel_abort = make(chan bool)

// Order channels
var orderChannel = make(chan formats.Order)
var floorChannel = make(chan formats.Floor)

func InitBackup() {

	fmt.Println("Backup routine has started.")
	// start go-rutiner for backup i nettverksmodulen
	go network.BackupCoordinator(backupChannel_read, backupChannel_write, channel_abort)
	// start go-rutiner for heartbeat

	// restore master (if necessary)
	go restoreMaster()
}

func InitMaster() {
	fmt.Println("Initializing master.")
	fmt.Println("IP address:", network.GetIP())
	// Control
	control.Init()
	// network
	go network.BackupWarden(backupChannel_read, backupChannel_write, channel_abort)
	go network.Warden(channel_read, channel_write, channel_abort)
	go network.Coordinator(channel_read, channel_write, channel_abort)
	// Backup
	generateBackup()
	// Update state machine
	stateMachine.SetMaster(true)
	// Events

	// Listener and broadcaster

	// Catch ctrl+c termination and stop the elevator

}

func generateBackup() {
	spawnCmd, err := exec.Command("gnome-terminal", "-x", "go", "run", "ex6.go")
	spawnCmd.Run()
	if err != nil {
		fmt.Println("Error: ", spawnCmd, err)
	}
	fmt.Println("A new backup has been spawned.")
}

func restoreMaster() {
	for {
		// Check if we are master
		if !stateMachine.IsMaster() {
			select {
			case <-channel_init_master:
				// We are no longer a backup
				initMaster()
				break
			}
		}
	}
}
