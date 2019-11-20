/*
	Implementation of the election algorithm by bullying to decide the
	coordinator process.
*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Since it is assumed that each process knows the ID of all the
// processes, the Process struct does not need to store the
// neighbours

type Process struct {
	id         int
	channel    chan Message
}

// Message.content is of `int` type
// 0 => OK
// 1 => ELECTION
// 2 => COORDINATOR
type Message struct {
	content int
	sender  Process
}

func processRunner(process Process, allProcesses []Process) {

	fmt.Printf("[Process %02d] Listening for messages\n", process.id)

	select {
	case message := <- process.channel:

		switch message.content {
		
		case 0:
			fmt.Printf("[Process %02d] Received OK from %v\n", process.id,
				message.sender.id)

			// Listen to other OK messages
			processRunner(process, allProcesses)
		
		case 1:
			fmt.Printf("[Process %02d] Received ELECTION from %v\n", process.id,
				message.sender.id)

			// If the process receives a ELECTION, send OK to the sender
			message.sender.channel <- Message{0, process}

			// Start election algorithm again
			go processRunner(process, allProcesses)
			startElection(process, allProcesses)
		
		case 2:
			fmt.Printf("[Process %02d] Received COORDINATOR from %v\n", process.id,
				message.sender.id)

			// Don't start the processRunner for this process anymore
		}
	
	case <- time.After(5 * time.Second):
		fmt.Printf("[Process %02d] Received no response for 5s\n", process.id)
		go processRunner(process, allProcesses)
	}
}

func startElection(process Process, allProcesses []Process) {
	fmt.Printf("[Process %02d] Starting election algorithm\n", process.id)

	var higherIDProcessExist bool = false

	for _, p := range allProcesses {
		if p.id > process.id {
			// Send a ELECTION message to Process p
			p.channel <- Message{1, process}

			higherIDProcessExist = true
		}
	}

	if !higherIDProcessExist {
		sendCordinatorMessage(process, allProcesses)
	}
}

func sendCordinatorMessage(process Process, allProcesses []Process) {
	for _, p := range allProcesses {
		// Send a COORDINATOR message to Process p
		p.channel <- Message{2, process}
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run tokenring.go <number of processes>")
		os.Exit(2)
	}

	numProcesses, _ := strconv.Atoi(os.Args[1])

	fmt.Println("Press any key to exit ...")
	fmt.Println("Starting program with", numProcesses, "processes")

	processes := make([]Process, numProcesses)

	for i := 0; i < numProcesses; i++ {
		processes[i] = Process{i + 1, make(chan Message)}
	}

	fmt.Println("\nCreated the following processes (with no cordinator):")
	fmt.Println("ID")
	fmt.Println("==")
	for i := 0; i < numProcesses; i++ {
		fmt.Printf("%02d\n", processes[i].id)
	}

	fmt.Println("")

	// Choose a random node to start the election algorithm
	rand.Seed(42)
	randomProcessIndex := rand.Intn(len(processes))
	randomProcess := processes[randomProcessIndex]

	// Start all processes
	for i := 0; i < numProcesses; i++ {
		go processRunner(processes[i], processes)
	}

	time.Sleep(time.Second * 1)
	startElection(randomProcess, processes)

	var input string
	fmt.Scanln(&input)
}