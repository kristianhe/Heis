package stateMachine

import (
	"fmt"
	"sync"
	//"../common" // Denne bare fjernes når jeg lagrer filen...?
)

var elevators []formats.Status
var mutex sync.Mutex

func AddExternalElevator(elevator formats.Status) {
	fmt.Println("Adding elevator", elevator.Elevator)
	mutex.Lock()
	defer mutex.Unlock()
	elevators = append(elevators, elevator)
}

func RemoveExternalElevator(elevator formats.Status) {
	localElevators := GetExternalElevators()
	// Check if we have elevators
	if len(localElevators) > 0 {
		for index := range localElevators {
			if localElevators[index].Elevator == elevator.Elevator {
				mutex.Lock()
				elevators = append(elevators[:index], elevators[index+1]...)
				mutex.Unlock()
				fmt.Println(elevator.Elevator, "is removed.")
				// To prevent out of bounds and panic we:
				break
			}
		}
	}
}

func UpdateExternalElevator(new formats.Status) {
	isFound := false
	localElevators := GetExternalElevators()
	// Check if we have elevators
	if len(localElevators) > 0 {
		for index := range localElevators {
			// Check if elevator exists
			if localElevators[index].Elevator == new.Elevator {
				mutex.Lock()
				elevators[index].State = new.State
				elevators[index].Floor = new.Floor
				elevators[index].Direction = new.Direction
				elevators[index].Time = new.Time
				mutex.Unlock()
				isFound = true
			}
		}
	}
	if !isFound {
		AddExternalElevator(new)
	}
}

func GetExternalElevators() []formats.Status {
	mutex.Lock()
	defer mutex.Unlock()
	// Create a copy to prevent data races
	copy := make([]formats.Status, len(elevators), len(elevators))
	// Need to manually copy all variables (the library function "copy" will not do)
	for id, elem := range elevators {
		copy[id] = elem
	}
	return copy
}

func CheckIfExternalElevatorExists(elevator formats.Status) bool {
	localElevators := GetExternalElevators()
	// Check if we have elevators
	if len(localElevators) > 0 {
		for index := range localElevators {
			if localElevators[index].Elevator == elevator.Elevator {
				return true
			}
		}
	}
}