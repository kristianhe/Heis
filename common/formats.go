// Standard message formats used in the network communication
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


//______________________________________________
// From KOK
//______________________________________________

/*

//Types 
type ID string

type Message struct {
Category 	int
Heartbeat 	Heartbeat
Status 		Status
Order		Order
Orders 		Orders
}

//Structs

type Heartbeat struct {
Counter int
}
type Status struct {
Elevator 	ID
State 		int
Floor 		int
Direction 	int
Priority	int
Time		time.Time
}
type Orders struct{
Elevator 	ID
List [] 	Order 
}
type Priority struct{
Elevator 	ID
Count		int
}
type Order struct {
Elevator 	ID 
Category	int 
Direction	int
Floor 		int  
Button      int
Time		time.Time
}
type Floor struct {
Current int
Status int
}
type Packet struct {
Address 	ID 
Data    	[]byte
}
*/
