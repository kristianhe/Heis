package spawn

import (
    "./network"

    "fmt"
)

//

func InitBackup() {                 // Initialization of backup
    fmt.Println("Backup routine has started.")
    network.ListenToClients()
}


func InitMaster() {                 // Initialization of master node

}

func generateBackup() {
    spawnCmd, err := exec.Command("gnome-terminal", "-x", "go", "run", "ex6.go")
	spawnCmd.Run()
    if err != nil {
        fmt.Println("Error: ", spawnCmd, err)
    }
    fmt.Println("A new backup has spawned.")
}
