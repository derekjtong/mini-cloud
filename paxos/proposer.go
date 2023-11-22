package paxos

import (
	"fmt"
	"net/rpc"
	"sync"
)

type Proposer struct {
	id             int
	proposalNumber int
	value          string
	acceptors      map[string]*rpc.Client // Given from node.go
	mu             sync.Mutex
}

func NewProposer(id int, proposalNumber int, acceptors map[string]*rpc.Client) *Proposer {
	return &Proposer{
		id:             id,
		proposalNumber: proposalNumber,
		acceptors:      acceptors,
	}
}

func (p *Proposer) Propose(value string) error {
	fmt.Printf("---PAXOS---\n")
	p.mu.Lock()
	p.value = value
	p.proposalNumber++
	p.mu.Unlock()

	fmt.Printf("---PHASE 1: PREPARE---\n")
	// Phase 1: Prepare
	prepareCount := 0
	for _, acceptor := range p.acceptors {
		response, err := p.sendPrepareRequest(acceptor, p.proposalNumber)
		if err != nil {
			fmt.Printf("Err")
			continue
		}
		if response.OK {
			prepareCount++
			if response.AcceptedValue != "" {
				p.mu.Lock()
				p.value = response.AcceptedValue
				p.mu.Unlock()
			}
		}
	}

	if prepareCount <= len(p.acceptors)/2 {
		return fmt.Errorf("failed to get majority in prepare phase")
	}

	fmt.Printf("accepted by %d nodes which is greater than the majority requirement of %d, moving to accept phase\n", prepareCount, len(p.acceptors)/2)

	fmt.Printf("---PHASE 2: ACCEPT---\n")
	// Phase 2: Accept
	acceptCount := 0
	for _, acceptor := range p.acceptors {
		response, err := p.sendAcceptRequest(acceptor, p.proposalNumber, p.value)
		if err != nil {
			continue
		}
		if response.OK {
			acceptCount++
		}
	}

	if acceptCount <= len(p.acceptors)/2 {
		return fmt.Errorf("failed to get majority in accept phase")
	}
	return nil
}

func (p *Proposer) sendPrepareRequest(acceptor *rpc.Client, proposalNumber int) (*PrepareResponse, error) {
	// TODO: Add timeout
	fmt.Printf("  node %d: send propose request\n", p.id)
	request := PrepareRequest{
		Id:       p.id,
		Proposal: proposalNumber,
	}
	var response PrepareResponse
	err := acceptor.Call("Node.Prepare", request, &response)

	fmt.Printf("  node %d: received prepare response from node %d\n", p.id, response.Id)
	return &response, err
}

func (p *Proposer) sendAcceptRequest(acceptor *rpc.Client, proposalNumber int, value string) (*AcceptResponse, error) {
	// TODO: Add timeout
	fmt.Printf("  node %d: sending accept request\n", p.id)
	request := AcceptRequest{
		Id:       p.id,
		Proposal: proposalNumber,
		Value:    value,
	}
	var response AcceptResponse
	err := acceptor.Call("Node.Accept", request, &response)
	fmt.Printf("  node %d: received accept response from node %d\n", p.id, response.Id)
	return &response, err
}
