package common

import (
	"time"
)

type ID string

type SimpleMessage struct {
	Address 	ID
	Data    	[]byte
}

type DetailedMessage struct {
	Category  	int
	Heartbeat 	Heartbeat
	Status    	Status
	Order     	Order
	Orders 		Orders
}

type Status struct {
	Elevator  	ID
	State     	int
	Floor     	int
	Direction 	int
	Priority  	int
	Time		time.Time
}

type Order struct {
	Category  	int
	Elevator  	ID
	Direction 	int
	Floor     	int
	Button    	int
	Time      	time.Time
}

type Orders struct {
	Elevator 	ID
	List     	[]Order
}

type Floor struct {
	Current 	int
	Status  	int
}

type Priority struct {
	Elevator 	ID
	Queue    	int
}

type Heartbeat struct {
	Count 		int
}
