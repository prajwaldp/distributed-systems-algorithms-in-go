package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Actor struct {
	id      int // ID of the actor
	counter int // The timestamp of the actor
}

type Message struct {
	from    Actor // The actor that sent the message
	to      Actor // The actor that receives the message
	message string // The message itself
	counter int // The logical clock counter
}

func send(sender Actor, receiver Actor, c chan Message) {
	var senderCounter int = sender.counter

	// Add one to the sender counter
	var updatedCounter int = senderCounter + 1

	msg := Message{sender, receiver, "Hello", updatedCounter}
	c <- msg
}

func receive(receiver Actor, c chan Message, actors []Actor) {
	for {
		msg := <- c
		
		receiver := msg.to
		receiverID := receiver.id
		receiverCounter := receiver.counter

		sender := msg.from
		senderID := sender.id
		senderCounter := sender.counter

		fmt.Println("[Received Message] | Sender ID:", senderID, "Counter:", senderCounter, "| Receiver ID:", receiverID, "Counter", receiverCounter)
		actors[receiverID - 1] = Actor{receiverID, receiverCounter + 1}
 	}
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run logicalclock.go num-actors num-messages")
		os.Exit(1)
	}

	numActors, err1 := strconv.Atoi(os.Args[1])
	numMessages, err2 := strconv.Atoi(os.Args[2])

	if err1 != nil || err2 != nil {
		fmt.Println("Usage: go run logicalclock.go num-actors num-messages")
		os.Exit(1)
	}

	fmt.Println("Initializing", numActors, "actors")
	actors := make([]Actor, numActors)

	for i := 0; i < numActors; i++ {
		actors[i] = Actor{i + 1, 0}
	}

	// Start the message passing queue (channel)
	var c chan Message = make(chan Message)

	for i := 0; i < numMessages; i++ {
		senderID := rand.Intn(numActors)
		receiverID := rand.Intn(numActors)

		for senderID == receiverID {
			senderID = rand.Intn(numActors)
			receiverID = rand.Intn(numActors)			
		}

		go send(actors[senderID], actors[receiverID], c)
		go receive(actors[receiverID], c, actors)

		time.Sleep(time.Second * 1)
	}

	fmt.Println("Press any key to exit")
	var input string
	fmt.Scanln(&input)
}
