package spawn

import (
	formats ".././common/formats"
	"../cases"
	"../control"
	"../network"
	"../orders"
	"../stateMachine"

	"fmt"
	"os/exec"
)

var filename string = "Spawn:"

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
var channel_order = make(chan formats.Order)
var channel_floor = make(chan formats.Floor)

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
	control.Init("localhost:15657")
	// Network goroutines
	go network.BackupWarden(backupChannel_read, backupChannel_write, channel_abort)
	go network.Warden(channel_read, channel_write, channel_abort)
	go network.Coordinator(channel_read, channel_write, channel_abort)
	// Backup
	generateBackup()
	// Update state machine
	stateMachine.SetMaster(true)
	// Case goroutines
	go cases.PollFloor(channel_floor)
	go cases.PollOrder(channel_order)
	go orders.Handle(channel_floor, channel_order, channel_write)
	// Listener and broadcaster
	go cases.Broadcaster(channel_write)
	go cases.ListenToNetwork(channel_read, channel_write)
	go cases.Heartbeater(backupChannel_write)
	go cases.SafetyCheck(channel_write)
	// Catch ctrl+c termination and stop the elevator
	go cases.ExitHandler()
}

func generateBackup() {
	spawnCmd := exec.Command("gnome-terminal", "-x", "go", "run", "main.go") // TODO Sjekk at main.go er riktig å skrive her
	spawnCmd.Run()
	fmt.Println("A new backup has been spawned.")
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
