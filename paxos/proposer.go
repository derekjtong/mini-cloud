package paxos

import (
	"fmt"
	"net/rpc"
)

type Proposer struct {
	id             int
	proposalNumber int
	value          string
	acceptors      map[string]*rpc.Client // Given from node.go

	highestAcceptedProposalNumber int
	highestAcceptedValue          string // Highest accepted value
}

func NewProposer(id int, proposalNumber int, acceptors map[string]*rpc.Client) *Proposer {
	return &Proposer{
		id:                            id,
		proposalNumber:                proposalNumber,
		acceptors:                     acceptors,
		highestAcceptedProposalNumber: -1,
	}
}

func (p *Proposer) Propose(value string) error {
	fmt.Printf("---PAXOS---\n")

	p.value = value
	proposalNumber := p.proposalNumber + 1
	p.proposalNumber = proposalNumber

	fmt.Printf("------PHASE 1: PREPARE------\n")
	// Phase 1: Prepare
	fmt.Printf("CURRENT PROPOSAL NUMBER: %d\n", p.proposalNumber)
	receivedPromises := 0

	for _, acceptor := range p.acceptors {
		response, err := p.sendPrepareRequest(acceptor, proposalNumber)
		if err != nil {
			continue
		}
		if response.OK {
			receivedPromises++
			if response.Proposal > p.highestAcceptedProposalNumber {
				fmt.Printf("%d > %d\n", response.Proposal, p.highestAcceptedProposalNumber)
				p.highestAcceptedProposalNumber = response.Proposal
				p.highestAcceptedValue = response.AcceptedValue
			}
		}
	}
	if receivedPromises <= len(p.acceptors)/2 {
		return fmt.Errorf("failed to get majority in prepare phase")
	}
	if p.highestAcceptedProposalNumber != -1 && p.highestAcceptedValue != "" {
		// Use the highest accepted value from the prepare phase
		p.value = p.highestAcceptedValue
	}

	fmt.Printf("accepted by %d nodes which is greater than the majority requirement of %d, moving to accept phase\n", receivedPromises, len(p.acceptors)/2)

	fmt.Printf("------PHASE 2: ACCEPT------\n")
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
	p.highestAcceptedProposalNumber = proposalNumber
	return nil
}

func (p *Proposer) sendPrepareRequest(acceptor *rpc.Client, proposalNumber int) (*PrepareResponse, error) {
	// TODO: Add timeout
	fmt.Printf("  node %d: send propose request {proposalNumber=%d}\n", p.id, proposalNumber)
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
	fmt.Printf("  node %d: sending accept request {proposalNumber=%d, value=%s}\n", p.id, proposalNumber, value)
	request := AcceptRequest{
		Id:       p.id,
		Proposal: proposalNumber,
		Value:    value,
	}
	var response AcceptResponse
	err := acceptor.Call("Node.Accept", request, &response)
	// fmt.Printf("  node %d: received accept response from node %d\n", p.id, response.Id)
	return &response, err
}
