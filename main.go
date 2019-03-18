package main

import (
	"flag"

	"./spawn"
	"./stateMachine"
)

func main() {

	flag_IsMaster := stateMachine.IsMaster()

	flag.BoolVar(&flag_IsMaster, "master", false, "X")
	flag.Parse()

	stateMachine.SetMaster(flag_IsMaster)

	if stateMachine.IsMaster() {

		spawn.InitMaster()

	} else if !stateMachine.IsMaster() {

		spawn.InitBackup()

	}

}
