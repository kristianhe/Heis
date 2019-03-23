package elevio

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var filename string = "Elevio: "

const POLL_RATE = 20 * time.Millisecond

var NUM_FLOORS int = 4

var _initialized bool = false
var mutex sync.Mutex
var conn net.Conn

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

// TODO Sette det ovenfor i en egen fil?

func Init(addr string, numFloors int) {

	if _initialized {

		fmt.Println("Driver already initialized.")

		return
	}
	NUM_FLOORS = numFloors
	mutex = sync.Mutex{}
	var err error
	conn, err = net.Dial("tcp", addr)
	if err != nil {

		panic(err.Error())

	}
	_initialized = true
}

func SetMotorDirection(dir MotorDirection) {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{1, byte(dir), 0, 0})

}

func SetButtonLamp(button ButtonType, floor int, value bool) {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{2, byte(button), byte(floor), toByte(value)})

}

func SetFloorIndicator(floor int) {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{3, byte(floor), 0, 0})

}

func SetDoorOpenLamp(value bool) {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{4, toByte(value), 0, 0})

}

func SetStopLamp(value bool) {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{5, toByte(value), 0, 0})

}

func PollButtons(receiver chan<- ButtonEvent) {

	prev := make([][3]bool, NUM_FLOORS)
	for {

		time.Sleep(POLL_RATE)
		for f := 0; f < NUM_FLOORS; f++ {

			for b := ButtonType(0); b < 3; b++ {

				v := getButton(b, f)
				if v != prev[f][b] && v != false {

					receiver <- ButtonEvent{f, ButtonType(b)}

				}
				prev[f][b] = v
			}
		}
	}

}

func PollFloorSensor(receiver chan<- int) {

	prev := -1
	for {
		time.Sleep(POLL_RATE)
		v := getFloor()
		if v != prev && v != -1 {
			receiver <- v
		}
		prev = v
	}

}

func PollStopButton(receiver chan<- bool) {

	prev := false
	for {
		time.Sleep(POLL_RATE)
		v := getStop()
		if v != prev {
			receiver <- v
		}
		prev = v
	}

}

func PollObstructionSwitch(receiver chan<- bool) {

	prev := false
	for {
		time.Sleep(POLL_RATE)
		v := getObstruction()
		if v != prev {
			receiver <- v
		}
		prev = v
	}

}

func getButton(button ButtonType, floor int) bool {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{6, byte(button), byte(floor), 0})
	var buf [4]byte
	conn.Read(buf[:])

	return toBool(buf[1])
}

func getFloor() int {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{7, 0, 0, 0})
	var buf [4]byte
	conn.Read(buf[:])
	if buf[1] != 0 {

		return int(buf[2])
	} else {

		return -1
	}

}

func getStop() bool {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{8, 0, 0, 0})
	var buf [4]byte
	conn.Read(buf[:])

	return toBool(buf[1])
}

func getObstruction() bool {

	mutex.Lock()
	defer mutex.Unlock()
	conn.Write([]byte{9, 0, 0, 0})
	var buf [4]byte
	conn.Read(buf[:])

	return toBool(buf[1])
}

func toByte(a bool) byte {

	var b byte = 0
	if a {
		b = 1
	}

	return b
}

func toBool(a byte) bool {

	var b bool = false
	if a != 0 {
		b = true
	}

	return b
}
