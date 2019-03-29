package spawn

import (
	formats "../common/formats"
	"../cases"
	"../control"
	"../network"
	"../orders"
	"../stateMachine"

	"fmt"
	"os/exec"
)

var filename string = "[Spawn] \t"

// Network channels
var channel_write 		= make(chan formats.SimpleMessage)
var channel_read 		= make(chan formats.SimpleMessage)
var backupChannel_write = make(chan formats.SimpleMessage)
var backupChannel_read 	= make(chan formats.SimpleMessage)

// Toggle channels
var channel_init_master = make(chan bool)
var channel_abort 		= make(chan bool)

// Order and floor channels
var channel_order = make(chan formats.Order)
var channel_floor = make(chan formats.Floor)

func InitBackup() {
	fmt.Println("Backup routine has started.")
	// Start goroutines for backup
	go network.BackupListener(backupChannel_read, channel_abort)
	// Start goroutines for heartbeat
	go cases.CheckHeartbeat(channel_abort, channel_init_master)
	go cases.CheckBackupHeartbeat(channel_init_master, backupChannel_read)
	// Restore master (if necessary)
	go restoreMaster()
}

func InitMaster() {
	fmt.Println("Initializing master.")
	fmt.Println("IP address:", network.GetIP())
	// Initialize control for hardware
	control.Init()
	// Start goroutines for network communication
	go network.BackupCoordinator(backupChannel_write, channel_abort)
	go network.MasterCoordinator(channel_read, channel_write, channel_abort)
	go network.SlaveCoordinator(channel_read, channel_write, channel_abort)
	// Spawn backup
	generateBackup()
	// Update state machine
	stateMachine.SetMaster(true)
	// Start goroutines for polling floors and orders
	go cases.PollFloor(channel_floor)
	go cases.PollOrder(channel_order)
	// Start goroutine for order handling
	go orders.Handle(channel_floor, channel_order, channel_write)
	// Start goroutines for broadcasting and listening
	go cases.Broadcaster(channel_write)
	go cases.ListenToNetwork(channel_read, channel_write)
	go cases.Heartbeater(backupChannel_write)
	go cases.SafeMode(channel_write)
	// Catch ctrl+c termination and stop the elevator
	go cases.ExitHandler()
}

func generateBackup() {
	// Spawn backup in new terminal
	spawnBackup := exec.Command("gnome-terminal", "-x", "go", "run", "main.go")
	spawnBackup.Run()
	fmt.Println("A new backup has been spawned.")
	// Spawn elevator server in new terminal
	spawnServer := exec.Command("gnome-terminal", "-x", "./ElevatorServer")
	spawnServer.Run()
	fmt.Println("A new elevator server has been spawned.")
}

func restoreMaster() {
	for {
		if !stateMachine.IsMaster() {
			select {
			case <-channel_init_master:
				// We are no longer a backup
				InitMaster()
				break
			}
		}
	}
}
