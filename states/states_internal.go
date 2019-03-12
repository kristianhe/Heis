package states

import (
	"../utilities"
	"fmt"
	"sync"
)

var mutex sync.mutex
var filename string = "State Machine -"
var connected_ip utilities.ID = "0.0.0.0"
var master bool = false
var state int = utilities.STATE_STARTUP
var direction int = utilities.INVALID
var floor int = 0
var connected int = utilities.DISCONNECTED
var priority int = 0


func IsMaster() bool (
	mutex.Lock
)