package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Actor struct {
	id               int
	channel          chan Token
	nextActorChannel chan Token // Each actor is connected to the next actor in a Ring Topology
}

type Token string

var x int = 0

func criticalSection(actorID int) {
	fmt.Printf("[Actor %d] Incremented x from %d\n", actorID, x)
	x = x + 1
	fmt.Printf("[Actor %d] Incremented x to %d\n", actorID, x)
}

func actorProcess(actor Actor, mutualExclusion bool) {
	
	// To simulate the fact that not all actors will need to execute the
	// critical section, only actors with even IDs will execute the
	// critical section and the others will not

	if actor.id % 2 == 0 {
		
		if mutualExclusion {
			// In a token-based mutual exclusion algorithm, the critical section
			// is only executed if the actor possesses the token

			token := <- actor.channel
			fmt.Printf("[Actor %d] Received the token \"%v\"\n", actor.id, token)

			if token == "token" {
				criticalSection(actor.id)
				fmt.Printf("[Actor %d] Forwarding the token to the next process\n", actor.id)
				actor.nextActorChannel <- token
			}
		
		} else {
			// Mutual exclusion is not enabled
			// The critical section is executed as is
			criticalSection(actor.id)
		}

	} else {

		if mutualExclusion {
			token := <- actor.channel
			fmt.Printf("[Actor %d] Received the token \"%v\"\n", actor.id, token)

			fmt.Printf("[Actor %d] Forwarding the token to the next process\n", actor.id)
			actor.nextActorChannel <- token
		}
	}

	// Run the actor process infinitely
	time.Sleep(5 * time.Second)
	fmt.Println("")
	actorProcess(actor, mutualExclusion)
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run tokenring.go [# actors] [enable/diable mutual exclusion --on or --off]")
		os.Exit(2)
	}

	numActors, _ := strconv.Atoi(os.Args[1])
	
	var mutualExclusion bool = true
	
	if os.Args[2] == "--off" {
		mutualExclusion = false
	}

	fmt.Println("Press any key to exit ...")
	fmt.Println("Starting program with", numActors, "actors")

	actors := make([]Actor, numActors)

	for i := 0; i < numActors; i++ {
		
		// Create an actor with a unique ID and a buffered channel of size 1
		// The buffered channel of size 1 is chosen because only one instance
		// of the token can exist in the entire network
		
		actors[i] = Actor{i + 1, make(chan Token, 1), nil}
	}

	for i := 0; i < numActors; i++ {
		actors[i].nextActorChannel = actors[(i + 1) % numActors].channel
	}

	fmt.Println("\nThe Ring Topology:\n")
	fmt.Println("Actor ID\t | Actor Channel\t | Connected Actor Channel")
	fmt.Println("------------------------------------------------------------------")
	for i := 0; i < numActors; i++ {
		actor := actors[i]
		fmt.Printf("%d\t\t | %v\t\t | %v\n", actor.id, actor.channel, actor.nextActorChannel)
	}

	fmt.Printf("\nStarting the actors...\n\n")

	// Starting all actors at once enabling/disabling Mutual Exclusion
	for i := 0; i < numActors; i++ {
		go actorProcess(actors[i], mutualExclusion)
	}

	actors[0].channel <- "token"

	var input string
	fmt.Scanln(&input)
}