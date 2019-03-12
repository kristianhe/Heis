package stateMachine

import (
	"sync"
)

var master bool = false
var mutex sync.Mutex

// Functions
func IsMaster() bool {

	mutex.Lock()
	defer mutex.Unlock()

	return master
}

func SetMaster(local_master bool) {

	mutex.Lock()
	defer mutex.Unlock()
	master = local_master

}
