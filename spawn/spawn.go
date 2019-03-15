package spawn

import (
	"../network"

	"fmt"
	"os/exec"
)

// Backup channels
var backupChannel_write = make(chan SimpleMessage)
var backupChannel_read = make(chan SimpleMessage)
var backupChannel_abort = make(chan bool)

// Master Channels
var channel_write = make(chan SimpleMessage)
var channel_read = make(chan SimpleMessage)
var orderChannel = make(chan Order)
var floorChannel = make(chan Floor)

func initBackup() {

	fmt.Println("Backup routine has started ...")
	// start go-rutiner for backup i nettverksmodulen
	go network.BackupCoordinator(backupChannel_read, backupChannel_write, backupChannel_abort)
	// start go-rutiner for heartbeat

	// restore master (if necessary)

	fmt.Println("... Backup routine is finished.")

}

func initMaster() {

}

func generateBackup() {

	spawnCmd, err := exec.Command("gnome-terminal", "-x", "go", "run", "ex6.go")
	spawnCmd.Run()
	if err != nil {
		fmt.Println("Error: ", spawnCmd, err)
	}
	fmt.Println("A new backup has been spawned.")

}
