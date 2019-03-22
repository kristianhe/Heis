package main

import (
	"./spawn"
	"./stateMachine"

	"fmt"
	"flag"
)

func main() {

	flag_isMaster := stateMachine.IsMaster()

	flag.BoolVar(&flag_isMaster, "master", false, "Start as master ??????")
	flag.Parse()

	stateMachine.SetMaster(flag_isMaster)

	if stateMachine.IsMaster() {
		fmt.Println("checkpoint1")
		spawn.InitMaster()

	} else if !stateMachine.IsMaster() {
		fmt.Println("checkpoint2")
		spawn.InitBackup()

	}

}
