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
	OrderList 	OrderList
}

type Status struct {
	Elevator  	ID
	State     	int
	Floor     	int
	Direction 	int
	Priority  	int // Bruke Priority-structen her?
	Time		time.Time
}

type Order struct {
	Category  	string
	Elevator  	ID
	Direction 	int
	Floor     	int
	Button    	int
	Time      	time.Time
}

type OrderList struct {
	Elevator 	ID
	List     	[]Order
}

type Floor struct {
	Current 	int
	Status  	int // Moving, idle etc. Bruke enum her i stedet for int?
}

type Priority struct {
	Elevator 	ID
	Queue    	int // Place in queue
}

type Heartbeat struct {
	Count 		int
}
