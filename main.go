package main

import (
	"./spawn"
	"./stateMachine"

	"flag"
	"time"
)

func main() {
	flag_isMaster := stateMachine.IsMaster()
	flag.BoolVar(&flag_isMaster, "master", false, "Start as master ??????")
	flag.Parse()
	stateMachine.SetMaster(flag_isMaster)

	if stateMachine.IsMaster() {
		spawn.InitMaster()
	} else if !stateMachine.IsMaster() {
		spawn.InitBackup()
	}

	// To prevent the system from stopping
	for {
		time.Sleep(time.Millisecond * 500)
	}
}
