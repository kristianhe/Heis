package stateMachine

import (
	formats ".././common/formats"

	"fmt"
)

var filename string = "[State Machine] "
var elevators []formats.Status

func addExternalElevator(elevator formats.Status) {
	fmt.Println("Adding elevator", elevator.Elevator)
	mutex.Lock()
	defer mutex.Unlock()
	elevators = append(elevators, elevator)
}

func RemoveExternalElevator(elevator formats.Status) {
	localElevators := GetExternalElevators()
	// Check if the list is emtpy
	if len(localElevators) > 0 {
		for index := range localElevators {
			if localElevators[index].Elevator == elevator.Elevator {
				mutex.Lock()
				elevators = elevators[:index+copy(elevators[index:], elevators[index+1:])]
				mutex.Unlock()
				fmt.Println(elevator.Elevator, "is removed.")
				break
			}
		}
	}
}

func UpdateExternalElevator(new formats.Status) {
	isFound := false
	localElevators := GetExternalElevators()
	// Check if the list is emtpy
	if len(localElevators) > 0 {
		for index := range localElevators {
			// Check if the elevator is in the list
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
		addExternalElevator(new)
	}
}

func GetExternalElevators() []formats.Status {
	mutex.Lock()
	defer mutex.Unlock()
	// Create a copy to prevent race conditions
	copy := make([]formats.Status, len(elevators), len(elevators))
	for id, elem := range elevators {
		copy[id] = elem
	}
	return copy
}
