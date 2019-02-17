package networkModule

import (
	"flag"
	"fmt"
	"os"
	"time"

	"./network/bcast"
	"./network/localip"
	"./network/peers"
)

// Message format (containing states etc.)
type Msg struct {
	Message string
	Iter    int
}

// TODO make a better name for this function
func NetworkFunc() {

	// Can set personal ID using two alternatives,
	// Alternative 1:
	// Choose by writing `go run main.go -id=your_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// Alternative 2:
	// Preset by local IP address + process ID
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// Channels for peer status
	peerUpdateCh := make(chan peers.PeerUpdate) // Peer status (who is currently active)
	peerTxEnable := make(chan bool)             // Enable/disable transmission

	// Start corresponding goroutines
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// Local channels for sending and receiving messages
	msgTx := make(chan Msg)
	msgRx := make(chan Msg)

	// Start broadcasting
	go bcast.Transmitter(16569, msgTx)
	go bcast.Receiver(16569, msgRx)

	// Send a message every second.
	go func() {
		helloMsg := Msg{"Hello from the one and only " + id, 0}
		for {
			helloMsg.Iter++
			msgTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	// While true
	fmt.Println()
	fmt.Println("----------------------")
	fmt.Println("NETWORK UP AND RUNNING")
	fmt.Println("----------------------")
	fmt.Println()
	for {
		select {
		case peerUpdate := <-peerUpdateCh:
			fmt.Println("----------------------------")
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peerUpdate.Peers)
			fmt.Printf("  New:      %q\n", peerUpdate.New)
			fmt.Printf("  Lost:     %q\n", peerUpdate.Lost)
			fmt.Println("----------------------------")

		case messageReceived := <-msgRx:
			fmt.Printf("Received: %#v\n", messageReceived)
		}
	}
}
