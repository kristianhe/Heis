package elevaluator

import (
	"time"

	".././common/constants"
	".././common/formats"
	"../control"
)

var heartbeat = time.Now()

func checkRequestedOrders() (int, int) {
	for floor := 0; floor < constants.FLOORS; floor++ {
		for button := 0; button < constants.BUTTONS; button++ {
			if control.GetButtonSignal(button, floor) == true {
				return button, floor
			}
		}
	}
}

// Vi har disse to i elevio
/*
func pollFloor(floorChannel chan formats.Floor) {
    return
}

func pollorder(orderChannel chan formats.Order) {
    return
}
*/

func statusChecker(channel_write chan formats.SimpleMessage) {
	return
}

func networkLoop(channel_write chan formats.SimpleMessage) {
	return
}

func networkListener(channel_read chan formats.SimpleMessage, channel_write chan formats.SimpleMessage) {
	return
}

func heartbeatLoop(backupChannel_write chan formats.SimpleMessage) {
	return
}

func heartbeatListener(channel_init_master chan bool, backupChannel_read chan formats.SimpleMessage) {
	return
}

func heartbeatChecker(channel_abort chan bool, channel_init_master chan bool) {
	return
}

// Starts a new heartbeat
func setHeartbeat() {
	heartbeat = time.Now()
}

func exitHandler() {
	return
}
