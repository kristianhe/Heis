package control
// Changes


import (
	"../elevio"
	"../states"
	"../utilities"
	"fmt"
)

var floor int
var filename string = "Controll -"

func Init() {

	//Initiate elevio
	var init_suc bool = elevio.Init()
	if init_suc != true {
		fmt.Println(filename, "Error when attempting to initialize ElevIO")
	}
	ClearLights()
	if GetFloorSignal() == utilities.INVALID {
		GoUp()
	}
	fmt.Println(filename, "Initialized ElevIO")
}

func GoUp() {
	DirUp()
	Move()
}

func GoDown() {
	DirDown()
	Move()
}

func DirUp() {
	elevio.ClearBit(elevio.MOTORDIR)
	states.SetDir(utilities.UP)
}





