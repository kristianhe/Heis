package control
// Changes


import (
	"../elevio"
	"../states"
	"../tools"
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
	if GetFloorSignal() == tools.INVALID {
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
	states.SetDir(tools.UP)
}

func DirDown() {
	elevio.SetBit(elevio.MOTORDIR)
	states.SetDir(tools.DOWN)
}

func DirectionSwitch() {
	if states.GetDir() == tools.UP {
		DirDown()
	} else {
		DirUp()
	}

}

func Move() {
	elevio.WriteAnalog(elevio.MOTOR, elevio.MOTOR_SPEED)
	states.SetState(tools.STATE_RUNNING)
}

func Stop() {
	elevio.WriteAnalog(elevio.MOTOR, tools.STOP)
}

func ClearLights() {
	for floor := 0; floor < tools.FLOORS; floor++ {
		for button := 0; button < tools.BUTTONS; button++ {
			SetButtonLamp(button, floor, tools.OFF)
		}
	}
	SetStopLamp(tools.OFF)
	SetDoorLamp(tools.OFF)
	SetFloorIndicator(tools.OFF)
}

func SetButtonLamp(button, floor, lamp int) int {

	if floor <= tools.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return tools.INVALID
	}
	if floor > tools.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", tools.FLOORS)
		return tools.INVALID
	}
	if button <= tools.INVALID {
		fmt.Println(filename, "Illegal button, must be larger than 0!")
		return tools.INVALID
	}
	if button > tools.BUTTONS {
		fmt.Println(filename, "Illegal button, must be less than ", tools.BUTTONS)
		return tools.INVALID
	}

	//Turn on lamp
	if lamp == tools.ON {
		elevio.SetBit(lamp_matrix[floor][button])
		return tools.TRUE
	}

	//Turn off lamp
	if lamp == tools.OFF {
		elevio.ClearBit(lamp_matrix[floor][button])
		return tools.TRUE
	}

	return tools.INVALID

}

func SetFloorIndicator(floor int) int {

	if floor <= tools.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return tools.INVALID
	}
	if floor > tools.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", tools.FLOORS)
		return tools.INVALID
	}

	switch floor {

	case tools.FLOOR_FIRST:
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND1)
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND2)
		return tools.TRUE
	case tools.FLOOR_SECOND:
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND1)
		elevio.SetBit(elevio.LIGHT_FLOOR_IND2)
		return tools.TRUE
	case tools.FLOOR_THIRD:
		elevio.SetBit(elevio.LIGHT_FLOOR_IND1)
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND2)
		return tools.TRUE
	case tools.FLOOR_LAST:
		elevio.SetBit(elevio.LIGHT_FLOOR_IND1)
		elevio.SetBit(elevio.LIGHT_FLOOR_IND2)
		return tools.TRUE

	}

	return tools.INVALID
}

func SetDoorLamp(lamp int) {
	if lamp == tools.ON {
		elevio.SetBit(elevio.LIGHT_DOOR_OPEN)
	}
	if lamp == tools.OFF {
		elevio.ClearBit(elevio.LIGHT_DOOR_OPEN)
	}
}

func SetStopLamp(lamp int) {
	if lamp == tools.ON {
		elevio.SetBit(elevio.LIGHT_STOP)
	}
	if lamp == tools.OFF {
		elevio.ClearBit(elevio.LIGHT_STOP)
	}
}

func GetButtonSignal(button, floor int) int {

	//Check if floor and button are valid
	if floor <= tools.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return tools.INVALID
	}
	if floor > tools.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", tools.FLOORS)
		return tools.INVALID
	}
	if button <= tools.INVALID {
		fmt.Println(filename, "Illegal button, must be larger than 0!")
		return tools.INVALID
	}
	if button > tools.BUTTONS {
		fmt.Println(filename, "Illegal button, must be less than ", tools.BUTTONS)
		return tools.INVALID
	}

	if elevio.ReadBit(button_matrix[floor][button]) == tools.TRUE {
		return tools.TRUE
	}

	return tools.FALSE
}

func GetFloorSignal() int {

	//Check all floors
	if elevio.ReadBit(elevio.SENSOR_FLOOR1) == tools.TRUE {
		return tools.FLOOR_FIRST
	}

	if elevio.ReadBit(elevio.SENSOR_FLOOR2) == tools.TRUE {
		return tools.FLOOR_SECOND
	}

	if elevio.ReadBit(elevio.SENSOR_FLOOR3) == tools.TRUE {
		return tools.FLOOR_THIRD
	}

	if elevio.ReadBit(elevio.SENSOR_FLOOR4) == tools.TRUE {
		return tools.FLOOR_LAST
	}

	//Invalid floor
	return tools.INVALID

}

//Returns true if the stop button is pressed
func GetStopSignal() int {
	return elevio.ReadBit(elevio.STOP)
}

//Returns true if we have a obstruction
func GetObstruction() int {
	return elevio.ReadBit(elevio.OBSTRUCTION)
}
