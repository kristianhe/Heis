package stateMachine

import (
	"../common"

	"sync"
)

var master bool = false
var mutex sync.Mutex
var connected int = constants.DISCONNECTED

func IsMaster() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return master
}

func IsConnected() bool {
	mutex.Lock()
	defer mutex.Unlock()
 	if connected == constants.CONNECTED { return true }
	return false
}

func SetMaster(local_master bool) {
	mutex.Lock()
	defer mutex.Unlock()
	master = local_master
}
