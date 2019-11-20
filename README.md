# Distributed systems algorithms in Golang

Simulation of distributed systems algorithms in Go.

Based on my coursework (COP5615 - Distributed Systems at the University of Florida),
the following algorithms are simulated:

1. Logical clocks (basic implementation of Lamport clocks)
2. Token ring algorithm (for safe execution of critical sections)
3. Election by bullying algorithm

## Caveats/Not Implemented

### Election by bullying

**Node Failure** - Suppose the process with ID 10 failed to respond, all
the go-routines will time out while waiting for the response. In this case,
the process with the ID 9, should be the coordinator as no-one responds to
its ELECTION request (while every other process has its ELECTION request
responded to by at least one OK response).
