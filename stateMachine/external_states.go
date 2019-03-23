package stateMachine

import (
	formats ".././common/formats"

	"fmt"
)

var filename string = "State Machine: "
var elevators []formats.Status

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
				elevators = elevators[:index+copy(elevators[index:], elevators[index+1:])]
				mutex.Unlock()
				fmt.Println(elevator.Elevator, "is removed.")
				// To prevent out of bounds and panic
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
	return false
}
