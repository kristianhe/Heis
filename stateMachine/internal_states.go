package stateMachine

import (
	constants ".././common/constants"
	formats ".././common/formats"

	"fmt"
	"sync"
)

var mutex sync.Mutex
var master bool = false
var connection int = constants.DISCONNECTED
var state int = constants.STATE_STARTUP
var floor int = 0
var direction int = constants.INVALID
var priority int = 0
var connectedIp formats.ID = "0.0.0.0"

func IsMaster() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return master
}

func IsConnected() bool {
	mutex.Lock()
	defer mutex.Unlock()
	if connection == constants.CONNECTED {
		return true
	}
	return false
}

func GetState() int {
	mutex.Lock()
	defer mutex.Unlock()
	return state
}

func GetFloor() int {
	mutex.Lock()
	defer mutex.Unlock()
	return floor
}

func GetDirection() int {
	mutex.Lock()
	defer mutex.Unlock()
	return direction
}

func GetPriority() int {
	mutex.Lock()
	defer mutex.Unlock()
	return priority
}

func GetConnectedIp() formats.ID {
	mutex.Lock()
	defer mutex.Unlock()
	return connectedIp
}

func PrintState() string {
	switch state {
	case constants.STATE_STARTUP:
		return "spawn"
	case constants.STATE_IDLE:
		return "idle"
	case constants.STATE_RUNNING:
		return "running"
	case constants.STATE_EMERGENCY:
		return "emergency"
	case constants.STATE_DOOR_OPEN:
		return "door open"
	case constants.STATE_DOOR_CLOSED:
		return "door closed"
	}
	return "invalid"
}

func PrintDirection() string {
	switch direction {
	case constants.UP:
		return "up"
	case constants.DOWN:
		return "down"
	case constants.STOP:
		return "stationary"
	}
	return "invalid"
}

func PrintFloor() string {
	switch floor {
	case constants.FLOOR_FIRST:
		return "1st"
	case constants.FLOOR_SECOND:
		return "2nd"
	case constants.FLOOR_THIRD:
		return "3rd"
	case constants.FLOOR_LAST:
		return "4th"
	}
	return "invalid"
}

func SetMaster(local_master bool) {
	mutex.Lock()
	defer mutex.Unlock()
	master = local_master
}

func SetConnection(desiredConnection int) {
	mutex.Lock()
	defer mutex.Unlock()
	if connection != desiredConnection {
		connection = desiredConnection
		fmt.Println(filename, "New connection ->", connection)
	}
}

func SetState(desiredState int) {
	mutex.Lock()
	defer mutex.Unlock()
	if state != desiredState {
		state = desiredState
		fmt.Println(filename, "New state ->", PrintState())
	}
}

func SetFloor(desiredFloor int) {
	mutex.Lock()
	defer mutex.Unlock()
	if floor != desiredFloor {
		floor = desiredFloor
		fmt.Println(filename, "New floor ->", PrintFloor())
	}
}

func SetDirection(desiredDirection int) {
	mutex.Lock()
	defer mutex.Unlock()
	if direction != desiredDirection {
		direction = desiredDirection
		fmt.Println(filename, "New direction ->", PrintDirection())
	}
}

func SetPriority(desiredPriority int) {
	mutex.Lock()
	defer mutex.Unlock()
	if priority != desiredPriority {
		priority = desiredPriority
		fmt.Println(filename, "New priority ->", priority)
	}
}

func SetConnectedIp(desiredConnectedIp formats.ID) {
	mutex.Lock()
	defer mutex.Unlock()
	if connectedIp != desiredConnectedIp {
		connectedIp = desiredConnectedIp
		fmt.Println(filename, "A new connection to", connectedIp, "has been established.")
	}
}
