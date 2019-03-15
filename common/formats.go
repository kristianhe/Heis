package common

import (
	"time"
)

type ID string

type SimpleMessage struct {
	Address ID
	Data    []byte
}

type DetailedMessage struct {
	Category  int // Message category ?? heller en enum?
	Heartbeat Heartbeat
	Status    Status
	Order     Order
	OrderList OrderList
}

type Status struct {
	Elevator  ID
	State     int
	Floor     int
	Direction int
	Priority  int // Bruke Priority-structen her?
}

type Order struct {
	Category  string // Message category ??
	Elevator  ID
	Direction int
	Floor     int
	Button    int
	time      time.Time
}

type OrderList struct {
	Elevator ID
	List     []Order
}

type Floor struct {
	Current int
	Status  int // Moving, idle etc. Bruke enum her i stedet for int?
}

type Priority struct {
	Elevator ID
	Queue    int // Place in queue
}

type Heartbeat struct {
	Count int
}
