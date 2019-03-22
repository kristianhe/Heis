package control

import (
	constants ".././common/constants"
	"../elevio"
	"../stateMachine"

	"fmt"
	"sync"
	"net"
)


var filename string = "Control -"
var mutex sync.Mutex
var conn net.Conn
var floor int
var initialized bool = false
var numFloors int = 4

func Init() {
	if initialized {
		fmt.Println("Driver already initialized!")
		return
	}
	var err error
	conn, err = net.Dial("tcp", "localhost:15657") // TODO Denne adressen må endres
	if err != nil	{ panic(err.Error()) }
	initialized = true
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
	conn.Write([]byte{1, byte(constants.UP), 0, 0})
	stateMachine.SetDirection(constants.UP)
}

func DirDown() {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{1, byte(constants.DOWN), 0, 0})
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
	conn.Write([]byte{1, byte(1), 0, 0})    		// 2800 is motor speed
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
	// Turn on lamp
	if lamp == constants.ON {
		mutex.Lock()
		defer mutex.Unlock()
		conn.Write([]byte{2, byte(button), byte(floor), toByte(true)})
		return constants.TRUE
	}
	// Turn off lamp
	if lamp == constants.OFF {
		mutex.Lock()
		defer mutex.Unlock()
		conn.Write([]byte{2, byte(button), byte(floor), toByte(false)})
		return constants.TRUE
	}
	return constants.INVALID
}

func SetFloorIndicator(floor int) int {
	conn.Write([]byte{3, byte(floor), 0, 0})
	/*
	switch floor {
	case constants.FLOOR_FIRST: //00
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND1)
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND2)

		return constants.TRUE
	case constants.FLOOR_SECOND: //01
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND1)
		elevio.SetBit(elevio.LIGHT_FLOOR_IND2)
		return constants.TRUE
	case constants.FLOOR_THIRD: //10
		elevio.SetBit(elevio.LIGHT_FLOOR_IND1)
		elevio.ClearBit(elevio.LIGHT_FLOOR_IND2)
		return constants.TRUE
	case constants.FLOOR_LAST: //11
		elevio.SetBit(elevio.LIGHT_FLOOR_IND1)
		elevio.SetBit(elevio.LIGHT_FLOOR_IND2)
		return constants.TRUE
	}
	*/
	return constants.INVALID
}

func SetDoorLamp(lamp int) {
	if lamp == constants.ON {
		mutex.Lock()
		defer mutex.Unlock()
		conn.Write([]byte{4, toByte(true), 0, 0})
	} else if lamp == constants.OFF {
		mutex.Lock()
		defer mutex.Unlock()
		conn.Write([]byte{4, toByte(false), 0, 0})
	}
}

func SetStopLamp(lamp int) {
	if lamp == constants.ON {
		mutex.Lock()
		defer mutex.Unlock()
		conn.Write([]byte{5, toByte(true), 0, 0})
	}
	if lamp == constants.OFF {
		mutex.Lock()
		defer mutex.Unlock()
		conn.Write([]byte{5, toByte(false), 0, 0})
	}
}

func GetButtonSignal(button, floor int) bool {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{6, byte(button), byte(floor), 0})
	var buffer [4]byte
	conn.Read(buffer[:])
	return toBool(buffer[1])
}

func GetFloorSignal() int {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{7, 0, 0, 0})
	var buffer [4]byte
	conn.Read(buffer[:])
	// Check all floors
	if buffer[1] != 0 {
		if int(buffer[2]) == elevio.SENSOR_FLOOR1 { return constants.FLOOR_FIRST }
		if int(buffer[2]) == elevio.SENSOR_FLOOR2 { return constants.FLOOR_SECOND }
		if int(buffer[2]) == elevio.SENSOR_FLOOR3 { return constants.FLOOR_THIRD }
		if int(buffer[2]) == elevio.SENSOR_FLOOR4 { return constants.FLOOR_LAST }
	}
	return constants.INVALID
}

// Returns 1 if the stop button is pressed
func GetStopSignal() int {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{8, 0, 0, 0})
	var buffer [4]byte
	conn.Read(buffer[:])
	return int(buffer[1])
}

// Returns 1 if we have a obstruction
func GetObstruction() int {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{9, 0, 0, 0})
	var buffer [4]byte
	conn.Read(buffer[:])
	return int(buffer[1])
}

func toBool(a byte) bool {
	var b bool = false
	if a != 0 	{ b = true }
	return b
}

func toByte(a bool) byte {
	var b byte = 0
	if a	{ b = 1	}
	return b
}
