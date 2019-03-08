package stateMachine

import (
    "fmt"
    "sync"
)

// Naming
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
    master = local_master
    defer mutex.Unlock()
}
