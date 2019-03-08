package stateMachine

import (
    "fmt"
    "sync"
)

var mutex sync.Mutex

func IsMaster() bool {
    mutex.Lock()
    defer mutex.Unlock()

    return false                        //Eventuelt en bool kalt master som er false
}
