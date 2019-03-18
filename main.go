package main

import (
	"./spawn"
	"./stateMachine"

	"flag"
)

func main() {

	flag_isMaster := stateMachine.IsMaster()

	flag.BoolVar(&flag_isMaster, "master", false, "X??")
	flag.Parse()

	stateMachine.SetMaster(flag_isMaster)

	if stateMachine.IsMaster() {

		spawn.InitMaster()

	} else if !stateMachine.IsMaster() {

		spawn.InitBackup()

	}

}
