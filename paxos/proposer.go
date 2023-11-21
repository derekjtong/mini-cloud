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
	acceptors      []*rpc.Client // Given from node.go
	mu             sync.Mutex
}

func NewProposer(id int, acceptors []*rpc.Client) *Proposer {
	return &Proposer{
		id:        id,
		acceptors: acceptors,
	}
}

func (p *Proposer) Propose(value string) error {
	p.mu.Lock()
	p.value = value
	p.proposalNumber++
	p.mu.Unlock()

	// Phase 1: Prepare
	prepareCount := 0
	for _, acceptor := range p.acceptors {
		response, err := p.sendPrepareRequest(acceptor, p.proposalNumber)
		if err != nil {
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
	request := PrepareRequest{
		Proposal: proposalNumber,
	}
	var response PrepareResponse
	err := acceptor.Call("Acceptor.Prepare", request, &response)
	return &response, err
}

func (p *Proposer) sendAcceptRequest(acceptor *rpc.Client, proposalNumber int, value string) (*AcceptResponse, error) {
	// TODO: Add timeout
	request := AcceptRequest{
		Proposal: proposalNumber,
		Value:    value,
	}
	var response AcceptResponse
	err := acceptor.Call("Acceptor.Accept", request, &response)
	return &response, err
}
