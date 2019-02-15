package orders

import(
	"time"
	
)

//Not sure where it should be in the folder structure
//Stud.ass: should create one transmitter with all required channels, not too big modules
func SendOrders(sender chan <- elevio.ButtonEvent){
    orderTx := make(chan elevio.ButtonEvent)
    go bcast.Transmitter(16569, orderTx)  
    
    go func() {
		for {
			orderTx <- sender
			time.Sleep(1 * time.Second)
		}
}