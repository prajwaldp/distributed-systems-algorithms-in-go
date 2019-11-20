package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Each actor maintains the timestamp (counter) in its state
type Actor struct {
	id      int
	counter int
	channel chan Msg
}

// Each message maintains the timestamp (counter) from the originting actor
type Msg struct {
	from    int
	to      int
	counter int
}

// Helper function to get the maximum timestamp from two timestamps
func getMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// Starts an actor
func startReceiving(actor Actor, actors []Actor) {
	fmt.Printf("[%d] Actor %d listening for messages\n", actor.counter, actor.id)
	
	// Read a message from the actors channel
	msg := <- actor.channel

	fmt.Printf("[%d] Actor %d got <Message from: %d, to: %d, counter: %d>\n",
		actor.counter, actor.id, msg.from, msg.to, msg.counter)

	// Choose the max of the two counters - the one with the actor and the one in the msg
	var updatedCounter int = getMax(actor.counter, msg.counter)

	// Increment counter by 1 [the event of processing the message and sending it]
	updatedCounter = updatedCounter + 1

	// Choose another random actor to send the message to
	randomActor := actors[rand.Intn(len(actors))]

	fmt.Printf("[%d] Actor %d sending to <Actor id: %d, counter: %d> <Message counter: %d>\n",
		updatedCounter, actor.id, randomActor.id, randomActor.counter, updatedCounter)
	
	// Send the message to the chosen random actors' channel
	randomActor.channel <- Msg{actor.id, randomActor.id, updatedCounter}

	// Update the actor struct
	updatedActor := Actor{actor.id, updatedCounter, actor.channel}

	// Start listening again
	time.Sleep(time.Second * 1)
	startReceiving(updatedActor, actors)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run logicalclock.go [actor count]")
		os.Exit(2)
	}

	numActors, _ := strconv.Atoi(os.Args[1])

	fmt.Println("Press any key to exit ...")
	fmt.Println("Starting program with", numActors, "actors")

	actors := make([]Actor, numActors)

	for i := 0; i < numActors; i++ {
		actor := Actor{i, 0, make(chan Msg)}
		actors[i] = actor
	}

	for i := 0; i < numActors; i++ {
		go startReceiving(actors[i], actors)
	}

	// Send a random actor a message
	rand.Seed(42)
	senderID := rand.Intn(len(actors))
	receiverID := rand.Intn(len(actors))

	receiver := actors[receiverID]

	// Start the message passing
	// The initial value of counter is 0
	receiver.channel <- Msg{senderID, receiverID, 0}

	var input string
	fmt.Scanln(&input)
}