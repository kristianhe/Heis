package control

import (
	constants ".././common/constants"
	"../stateMachine"

	"fmt"
	"net"
	"sync"
)

var filename string = "Control: "
var mutex sync.Mutex
var conn net.Conn
var floor int
var isInitialized bool = false
var NUM_FLOORS int = 4

func Init(addr string) {
	if isInitialized {
		fmt.Println("Driver already initialized!")
		return
	}
	//ClearLights()
	var err error
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	isInitialized = true
}

/*
func GoUp() {
	DirUp()
	// Move()
}

func GoDown() {
	DirDown()
	// Move()
}
*/

func DirUp() {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{1, byte(1), 0, 0})
	stateMachine.SetDirection(constants.UP)
	stateMachine.SetState(constants.STATE_RUNNING)
}

func DirDown() {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{1, byte(1), 0, 0}) // TODO SKIFT FRA 1 til -1
	stateMachine.SetDirection(constants.DOWN)
	stateMachine.SetState(constants.STATE_RUNNING)
}

func SwitchDir() { // TODO denne er det noe muffins med...
	if stateMachine.GetDirection() == constants.UP {
		DirDown()
	} else {
		DirUp()
	}
}

// func Move() {
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	conn.Write([]byte{1, byte(2800), 0, 0}) // 2800 is motor speed
// 	stateMachine.SetState(constants.STATE_RUNNING)
// }

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

func SetFloorIndicator(floor int) {
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{3, byte(floor), 0, 0})
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

func GetFloorSignal() int { // TODO Denne funksjonen må endres, trenger ikke if
	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{7, 0, 0, 0})
	var buffer [4]byte
	conn.Read(buffer[:])
	// Check all floors
	if buffer[1] != 0 {
		if int(buffer[2]) == constants.FLOOR_FIRST {
			return constants.FLOOR_FIRST
		}
		if int(buffer[2]) == constants.FLOOR_SECOND {
			return constants.FLOOR_SECOND
		}
		if int(buffer[2]) == constants.FLOOR_THIRD {
			return constants.FLOOR_THIRD
		}
		if int(buffer[2]) == constants.FLOOR_LAST {
			return constants.FLOOR_LAST
		}
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
	if a != 0 {
		b = true
	}
	return b
}

func toByte(a bool) byte {
	var b byte = 0
	if a {
		b = 1
	}
	return b
}
