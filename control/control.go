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

func DirDown() {
	elevio.SetBit(elevio.MOTORDIR)
	states.SetDir(utilities.DOWN)
}

func DirectionSwitch() {
	if states.GetDir() == utilities.UP {
		DirDown()
	} else {
		DirUp()
	}

}

func Move() {
	elevio.WriteAnalog(elevio.MOTOR, elevio.MOTOR_SPEED)
	states.SetState(utilities.STATE_RUNNING)
}

func Stop() {
	elevio.WriteAnalog(elevio.MOTOR, utilities.STOP)
}

func ClearLights() {
	for floor := 0; floor < utilities.FLOORS; floor++ {
		for button := 0; button < utilities.BUTTONS; button++ {
			SetButtonLamp(button, floor, utilities.OFF)
		}
	}
	SetStopLamp(utilities.OFF)
	SetDoorLamp(utilities.OFF)
	SetFloorIndicator(utilities.OFF)
}

func SetButtonLamp(button, floor, lamp int) int {

	if floor <= utilities.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return utilities.INVALID
	}
	if floor > utilities.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", utilities.FLOORS)
		return utilities.INVALID
	}
	if button <= utilities.INVALID {
		fmt.Println(filename, "Illegal button, must be larger than 0!")
		return utilities.INVALID
	}
	if button > utilities.BUTTONS {
		fmt.Println(filename, "Illegal button, must be less than ", utilities.BUTTONS)
		return utilities.INVALID
	}

	//Turn on lamp
	if lamp == utilities.ON {
		hardware.SetBit(lamp_matrix[floor][button])
		return utilities.TRUE
	}

	//Turn off lamp
	if lamp == utilities.OFF {
		hardware.ClearBit(lamp_matrix[floor][button])
		return utilities.TRUE
	}

	return utilities.INVALID

}

func SetFloorIndicator(floor int) int {

	if floor <= utilities.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return utilities.INVALID
	}
	if floor > utilities.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", utilities.FLOORS)
		return utilities.INVALID
	}

	switch floor {

	case utilities.FLOOR_FIRST:
		hardware.ClearBit(hardware.LIGHT_FLOOR_IND1)
		hardware.ClearBit(hardware.LIGHT_FLOOR_IND2)
		return utilities.TRUE
	case utilities.FLOOR_SECOND:
		hardware.ClearBit(hardware.LIGHT_FLOOR_IND1)
		hardware.SetBit(hardware.LIGHT_FLOOR_IND2)
		return utilities.TRUE
	case utilities.FLOOR_THIRD:
		hardware.SetBit(hardware.LIGHT_FLOOR_IND1)
		hardware.ClearBit(hardware.LIGHT_FLOOR_IND2)
		return utilities.TRUE
	case utilities.FLOOR_LAST:
		hardware.SetBit(hardware.LIGHT_FLOOR_IND1)
		hardware.SetBit(hardware.LIGHT_FLOOR_IND2)
		return utilities.TRUE

	}

	return utilities.INVALID
}

func SetDoorLamp(lamp int) {
	if lamp == utilities.ON {
		hardware.SetBit(hardware.LIGHT_DOOR_OPEN)
	}
	if lamp == utilities.OFF {
		hardware.ClearBit(hardware.LIGHT_DOOR_OPEN)
	}
}

func SetStopLamp(lamp int) {
	if lamp == utilities.ON {
		hardware.SetBit(hardware.LIGHT_STOP)
	}
	if lamp == utilities.OFF {
		hardware.ClearBit(hardware.LIGHT_STOP)
	}
}

func GetButtonSignal(button, floor int) int {

	//Check if floor and button are valid
	if floor <= utilities.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return utilities.INVALID
	}
	if floor > utilities.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", utilities.FLOORS)
		return utilities.INVALID
	}
	if button <= utilities.INVALID {
		fmt.Println(filename, "Illegal button, must be larger than 0!")
		return utilities.INVALID
	}
	if button > utilities.BUTTONS {
		fmt.Println(filename, "Illegal button, must be less than ", utilities.BUTTONS)
		return utilities.INVALID
	}

	if hardware.ReadBit(button_matrix[floor][button]) == utilities.TRUE {
		return utilities.TRUE
	}

	return utilities.FALSE
}

func GetFloorSignal() int {

	//Check all floors
	if hardware.ReadBit(hardware.SENSOR_FLOOR1) == utilities.TRUE {
		return utilities.FLOOR_FIRST
	}

	if hardware.ReadBit(hardware.SENSOR_FLOOR2) == utilities.TRUE {
		return utilities.FLOOR_SECOND
	}

	if hardware.ReadBit(hardware.SENSOR_FLOOR3) == utilities.TRUE {
		return utilities.FLOOR_THIRD
	}

	if hardware.ReadBit(hardware.SENSOR_FLOOR4) == utilities.TRUE {
		return utilities.FLOOR_LAST
	}

	//Invalid floor
	return utilities.INVALID

}

//Returns true if the stop button is pressed
func GetStopSignal() int {
	return hardware.ReadBit(hardware.STOP)
}

//Returns true if we have a obstruction
func GetObstruction() int {
	return hardware.ReadBit(hardware.OBSTRUCTION)
}