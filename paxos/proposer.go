package paxos

import (
	"fmt"
	"net/rpc"
)

type Proposer struct {
	id             int
	ProposalNumber int
	Value          string
	Acceptors      map[string]*rpc.Client // Given from node.go

	HighestAcceptedProposalNumber int
	HighestAcceptedValue          string // Highest accepted value
}

func NewProposer(id int, proposalNumber int, acceptors map[string]*rpc.Client) *Proposer {
	return &Proposer{
		id:                            id,
		ProposalNumber:                proposalNumber,
		Acceptors:                     acceptors,
		HighestAcceptedProposalNumber: -1,
	}
}

func (p *Proposer) Propose(value string) error {
	fmt.Printf("---PAXOS---\n")

	p.Value = value
	proposalNumber := p.ProposalNumber + 1
	p.ProposalNumber = proposalNumber

	fmt.Printf("------PHASE 1: PREPARE------\n")
	// Phase 1: Prepare
	fmt.Printf("Current proposal number: %d\n", p.ProposalNumber)
	receivedPromises := 0
	// p.HighestAcceptedProposalNumber = -1
	p.HighestAcceptedValue = ""
	for _, acceptor := range p.Acceptors {
		response, err := p.sendPrepareRequest(acceptor, proposalNumber)
		if err != nil {
			continue
		}
		if response.OK {
			receivedPromises++
			if response.Proposal > p.HighestAcceptedProposalNumber {
				fmt.Printf("!!HIGHER ACCEPTED PROPOSAL DETECTED (%d > %d) - Changing accepted value from '%s' to '%s'\n", response.Proposal, p.HighestAcceptedProposalNumber, p.HighestAcceptedValue, response.AcceptedValue)
				p.HighestAcceptedProposalNumber = response.Proposal
				p.HighestAcceptedValue = response.AcceptedValue
			}
		}
	}
	if p.HighestAcceptedProposalNumber != -1 && p.HighestAcceptedValue != "" {
		// Use the highest accepted value from the prepare phase
		fmt.Printf("!!SENDING NEW VALUE - Changing send value from '%s' to '%s'\n", p.Value, p.HighestAcceptedValue)
		p.Value = p.HighestAcceptedValue
	}

	if receivedPromises < len(p.Acceptors)/2+1 {
		fmt.Printf("failed to gain consensus: accepted by %d nodes which is less than the majority requirement of %d\n", receivedPromises, len(p.Acceptors)/2+1)
		return fmt.Errorf("failed to get majority in prepare phase")
	}
	fmt.Printf("proceeding to accept phase: accepted by %d nodes which is greater than the majority requirement of %d\n", receivedPromises, len(p.Acceptors)/2+1)

	fmt.Printf("------PHASE 2: ACCEPT------\n")
	fmt.Printf("Sending %s\n", p.Value)
	acceptCount := 0
	for _, acceptor := range p.Acceptors {
		response, err := p.sendAcceptRequest(acceptor, p.ProposalNumber, p.Value)
		if err != nil {
			continue
		}
		if response.OK {
			acceptCount++
		}
	}

	if acceptCount < len(p.Acceptors)/2+1 {
		return fmt.Errorf("failed to get majority in accept phase")
	}
	p.HighestAcceptedProposalNumber = proposalNumber
	return nil
}

func (p *Proposer) sendPrepareRequest(acceptor *rpc.Client, proposalNumber int) (*PrepareResponse, error) {
	// TODO: Add timeout
	request := PrepareRequest{
		Id:       p.id,
		Proposal: proposalNumber,
	}
	fmt.Printf("  node %d: sending  - %#v\n", p.id, request)
	var response PrepareResponse
	err := acceptor.Call("Node.Prepare", request, &response)

	fmt.Printf("  node %d: received - %#v\n", p.id, response)
	return &response, err
}

func (p *Proposer) sendAcceptRequest(acceptor *rpc.Client, proposalNumber int, value string) (*AcceptResponse, error) {
	// TODO: Add timeout
	request := AcceptRequest{
		Id:       p.id,
		Proposal: proposalNumber,
		Value:    value,
	}
	fmt.Printf("  node %d: sending  - %#v\n", p.id, request)
	var response AcceptResponse
	err := acceptor.Call("Node.Accept", request, &response)

	// fmt.Printf("  node %d: received %#v\n", p.id, response)
	return &response, err
}
