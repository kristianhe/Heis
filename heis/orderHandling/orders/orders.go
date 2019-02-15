package orders

import(
	"time"
	"./../elevio"
	"./../networkModule"
)

func SendOrders(sender chan <- elevio.ButtonEvent){
    orderTx := make(chan elevio.ButtonEvent)
    go bcast.Transmitter(16569, orderTx)  
    
    go func() {
		for {
			orderTx <- sender
			time.Sleep(1 * time.Second)
		}
}