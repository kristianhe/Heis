package control

import (
	".././common/constants"
	"../elevio"
	"../stateMachine"

	"fmt"
)

var floor int
var filename string = "Control -"

func Init() {
	// Initiate elevio
	var init_suc bool = elevio.Init()
	if init_suc != true {
		fmt.Println(filename, "Error when attempting to initialize ElevIO")
	}
	ClearLights()
	if GetFloorSignal() == constants.INVALID {
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
	stateMachine.SetDir(constants.UP)
}

func DirDown() {
	elevio.SetBit(elevio.MOTORDIR)
	stateMachine.SetDir(constants.DOWN)
}

func SwitchDir() { // TODO skiftet fra "DirectionSwitch" til SwitchDir
	if stateMachine.GetDir() == constants.UP {
		DirDown()
	} else {
		DirUp()
	}
}

func Move() {
	elevio.WriteAnalog(elevio.MOTOR, elevio.MOTOR_SPEED)
	stateMachine.SetState(constants.STATE_RUNNING)
}

func Stop() {
	elevio.WriteAnalog(elevio.MOTOR, constants.STOP)
}

func ClearLights() {
	for floor := 0; floor < constants.FLOORS; floor++ {
		for button := 0; button < constants.BUTTONS; button++ {
			SetButtonLamp(button, floor, constants.OFF)
		}
	}
	SetStopLamp(constants.OFF)
	SetDoorLamp(constants.OFF)
	SetFloorIndicator(constants.OFF)
}

func SetButtonLamp(button, floor, lamp int) int {
	if floor <= constants.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return constants.INVALID
	}
	if floor > constants.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", constants.FLOORS)
		return constants.INVALID
	}
	if button <= constants.INVALID {
		fmt.Println(filename, "Illegal button, must be larger than 0!")
		return constants.INVALID
	}
	if button > constants.BUTTONS {
		fmt.Println(filename, "Illegal button, must be less than ", constants.BUTTONS)
		return constants.INVALID
	}
	// Turn on lamp
	if lamp == constants.ON {
		elevio.SetBit(lamp_matrix[floor][button])
		return constants.TRUE
	}
	// Turn off lamp
	if lamp == constants.OFF {
		elevio.ClearBit(lamp_matrix[floor][button])
		return constants.TRUE
	}
	return constants.INVALID
}

func SetFloorIndicator(floor int) int {
	if floor <= constants.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return constants.INVALID
	}
	if floor > constants.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", constants.FLOORS)
		return constants.INVALID
	}
	switch floor {
	case constants.FLOOR_FIRST:
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND1)
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND2)
		return constants.TRUE
	case constants.FLOOR_SECOND:
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND1)
		elevio.SetBit(elevio.LIGHT_FLOOR_IND2)
		return constants.TRUE
	case constants.FLOOR_THIRD:
		elevio.SetBit(elevio.LIGHT_FLOOR_IND1)
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND2)
		return constants.TRUE
	case constants.FLOOR_LAST:
		elevio.SetBit(elevio.LIGHT_FLOOR_IND1)
		elevio.SetBit(elevio.LIGHT_FLOOR_IND2)
		return constants.TRUE
	}
	return constants.INVALID
}

func SetDoorLamp(lamp int) {
	if lamp == constants.ON {
		elevio.SetBit(elevio.LIGHT_DOOR_OPEN)
	}
	if lamp == constants.OFF {
		elevio.ClearBit(elevio.LIGHT_DOOR_OPEN)
	}
}

func SetStopLamp(lamp int) {
	if lamp == constants.ON {
		elevio.SetBit(elevio.LIGHT_STOP)
	}
	if lamp == constants.OFF {
		elevio.ClearBit(elevio.LIGHT_STOP)
	}
}

func GetButtonSignal(button, floor int) int {
	// Check if floor and button are valid
	if floor <= constants.INVALID {
		fmt.Println(filename, "Illegal floor, must be larger than 0!")
		return constants.INVALID
	}
	if floor > constants.FLOORS {
		fmt.Println(filename, "Illegal floor, must be less than ", constants.FLOORS)
		return constants.INVALID
	}
	if button <= constants.INVALID {
		fmt.Println(filename, "Illegal button, must be larger than 0!")
		return constants.INVALID
	}
	if button > constants.BUTTONS {
		fmt.Println(filename, "Illegal button, must be less than ", constants.BUTTONS)
		return constants.INVALID
	}
	if elevio.ReadBit(button_matrix[floor][button]) == constants.TRUE {
		return constants.TRUE
	} else {
		return constants.FALSE
	}
}

func GetFloorSignal() int {
	// Check all floors
	if elevio.ReadBit(elevio.SENSOR_FLOOR1) == constants.TRUE {
		return constants.FLOOR_FIRST
	}
	if elevio.ReadBit(elevio.SENSOR_FLOOR2) == constants.TRUE {
		return constants.FLOOR_SECOND
	}
	if elevio.ReadBit(elevio.SENSOR_FLOOR3) == constants.TRUE {
		return constants.FLOOR_THIRD
	}
	if elevio.ReadBit(elevio.SENSOR_FLOOR4) == constants.TRUE {
		return constants.FLOOR_LAST
	}
	//Invalid floor
	return constants.INVALID
}

// Returns true if the stop button is pressed
func GetStopSignal() int { return elevio.ReadBit(elevio.STOP) }

// Returns true if we have a obstruction
func GetObstruction() int { return elevio.ReadBit(elevio.OBSTRUCTION) }
