package stateMachine

import (
constants	".././common/constants"
formats	".././common/formats"

	"fmt"
	"sync"
)

var mute sync.Mutex
var master bool = false
var connection int = constants.DISCONNECTED
var state int = constants.STATE_STARTUP
var floor int = 0
var direction int = constants.INVALID
var priority int = 0
var connectedIp formats.ID = "0.0.0.0"

func IsMaster() bool {
	mute.Lock()
	defer mute.Unlock()
	return master
}

func IsConnected() bool {
	mute.Lock()
	defer mute.Unlock()
	if connection == constants.CONNECTED  { return true }
	return false
}

func GetState() int {
	mute.Lock()
	defer mute.Unlock()
	return state
}

func GetFloor() int {
	mute.Lock()
	defer mute.Unlock()
	return floor
}

func GetDirection() int {
	mute.Lock()
	defer mute.Unlock()
	return direction
}

func GetPriority() int {
	mute.Lock()
	defer mute.Unlock()
	return priority
}

func GetConnectedIp() formats.ID {
	mute.Lock()
	defer mute.Unlock()
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

func SetMaster(local_master bool) {
	mute.Lock()
	defer mute.Unlock()
	master = local_master
}

func SetConnection(desiredConnection int) {
	mute.Lock()
	defer mute.Unlock()
	if connection != desiredConnection {
		connection = desiredConnection
		fmt.Println("Setting connection to", connection)
	}
}

func SetState(desiredState int) {
	mute.Lock()
	defer mute.Unlock()
	if state != desiredState {
		state = desiredState
		fmt.Println(state, "is the new state.")
	}
}

func SetFloor(desiredFloor int) {
	mute.Lock()
	defer mute.Unlock()
	if floor != desiredFloor {
		floor = desiredFloor
		fmt.Println("Setting floor to", floor)
	}
}

func SetDirection(desiredDirection int) {
	mute.Lock()
	defer mute.Unlock()
	if direction != desiredDirection {
		direction = desiredDirection
		fmt.Println("Setting direction to", direction)
	}
}

func SetPriority(desiredPriority int) {
	mute.Lock()
	defer mute.Unlock()
	if priority != desiredPriority {
		priority = desiredPriority
		fmt.Println(priority, "is the new priority.")
	}
}

func SetConnectedIp(desiredConnectedIp formats.ID) {
	mute.Lock()
	defer mute.Unlock()
	if connectedIp != desiredConnectedIp {
		connectedIp = desiredConnectedIp
		fmt.Println("A new connection to", connectedIp, "has been established.")
	}
}
