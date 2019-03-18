package spawn

import (
	".././common/formats"
	"../control"
	"../network"
	"../cases"
	"../orders"

	"fmt"
	"os/exec"
)

// Main channels											// TODO revurder navn på channels... vær konsekvent
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
	// Start goroutines for backup
	go network.BackupCoordinator(backupChannel_read, backupChannel_write, channel_abort)
	// Start goroutines for heartbeat
	go cases.CheckHeartbeat(channel_abort, channel_init_master)
	go cases.CheckBackupHeartbeat(channel_init_master, backupChannel_read)
	// Restore master (if necessary)
	go restoreMaster()
}

func InitMaster() {
	fmt.Println("Initializing master.")
	fmt.Println("IP address:", network.GetIP())
	// Control
	control.Init()
	// Network goroutines
	go network.BackupWarden(backupChannel_read, backupChannel_write, channel_abort)
	go network.Warden(channel_read, channel_write, channel_abort)
	go network.Coordinator(channel_read, channel_write, channel_abort)
	// Backup
	generateBackup()
	// Update state machine
	stateMachine.SetMaster(true)
	// Case goroutines
	go cases.PollFloor(floorChannel)
	go cases.PollOrder(orderChannel)
	go orders.Handle(floorChannel, orderChannel, channel_write)
	// Listener and broadcaster
	go cases.BroadcastToNetwork(channel_write)
	go cases.ListenToNetwork(channel_read, channel_write)
	go cases.Heartbeater(backupChannel_write)
	go cases.CheckStatus(channel_write)
	// Catch ctrl+c termination and stop the elevator
	go cases.ExitHandler()
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
