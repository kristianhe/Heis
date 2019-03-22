package control

import (
	".././common/constants"
	"../elevio"
	"../stateMachine"

	"fmt"
	"net"
)

var floor int
var filename string = "Control -"

func Init() {
	// Initiate elevio
	var initSuccess bool = elevio.Init()
	if initSuccess != true   { fmt.Println(filename, "Error when attempting to initialize ElevIO") }
	ClearLights()
	if GetFloorSignal() == constants.INVALID   { GoUp() }
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
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{1, byte(constants.MOTOR_UP), 0, 0})
	stateMachine.SetDirection(constants.UP)
}

func DirDown() {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{1, byte(constants.MOTOR_DOWN), 0, 0})
	stateMachine.SetDirection(constants.DOWN)
}

func SwitchDir() { 										// TODO skiftet fra "DirectionSwitch" til SwitchDir
	if stateMachine.GetDirection() == constants.UP {
		DirDown()
	} else {
		DirUp()
	}
}

func Move() {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{1, byte(2800), 0, 0})    		// 2800 is motor speed
	stateMachine.SetState(constants.STATE_RUNNING)
}

func Stop() {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{1, byte(constants.STOP), 0, 0})
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
	} else if lamp == constants.OFF {
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
	// Invalid floor
	return constants.INVALID
}

// Returns true if the stop button is pressed
func GetStopSignal() int   { return elevio.ReadBit(elevio.STOP) }

// Returns true if we have a obstruction
func GetObstruction() int   { return elevio.ReadBit(elevio.OBSTRUCTION) }
