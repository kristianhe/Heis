package stateMachine

import (
	"sync"
)

var master bool = false
var mutex sync.Mutex

func IsMaster() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return master
}

func IsConnected() bool {
	mutex.Lock()
	defer mutex.Unlock()
 
}

func SetMaster(local_master bool) {
	mutex.Lock()
	defer mutex.Unlock()
	master = local_master
}
