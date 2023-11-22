package paxos

import (
	"fmt"
)

type Acceptor struct {
	Id               int
	PromisedProposal int    // Highest prepare request seen so far
	AcceptedProposal int    // Highest proposal agreed upon
	AcceptedValue    string // Value of the highest proposal agreed upon
}

func NewAcceptor(id int) *Acceptor {
	return &Acceptor{
		Id:               id,
		PromisedProposal: -1,
		AcceptedProposal: -1,
		AcceptedValue:    "",
	}
}

// Handle Prepare request
func (a *Acceptor) Prepare(proposal int) PrepareResponse {
	fmt.Printf("    node %d: STATUS   - %#v\n", a.Id, a)
	if proposal > a.PromisedProposal {
		fmt.Printf("    node %d: ACCEPTED - changing PromisedProposal from %d to incoming %d\n", a.Id, a.PromisedProposal, proposal)
		a.PromisedProposal = proposal
		fmt.Printf("    node %d: STATUS   - %#v\n", a.Id, a)
		// fmt.Printf("    node %d: {promisedProposal=%d, acceptedValue=%s, acceptedProposal=%d}\n", a.Id, a.PromisedProposal, a.AcceptedValue, a.AcceptedProposal)
		// Promise to not accept any earlier proposals
		return PrepareResponse{
			Id:            a.Id,
			OK:            true,
			Proposal:      a.AcceptedProposal,
			AcceptedValue: a.AcceptedValue,
		}
	}
	fmt.Printf("    node %d: REJECTED - PromisedProposal:%d greater than incoming proposal: %d\n", a.Id, a.PromisedProposal, proposal)
	return PrepareResponse{
		Id: a.Id,
		OK: false,
	}
}

// Handle Accept request
func (a *Acceptor) Accept(proposal int, value string) AcceptResponse {
	fmt.Printf("    node %d: STATUS   - %#v\n", a.Id, a)

	if proposal >= a.PromisedProposal {
		fmt.Printf("    node %d: ACCEPTED - proposal:%d >= PromisedProposal:%d\n", a.Id, proposal, a.PromisedProposal)
		fmt.Printf("    node %d: UPDATING - value '%s' -> '%s', AcceptedProposal %d -> %d\n", a.Id, a.AcceptedValue, value, a.AcceptedProposal, proposal)
		a.PromisedProposal = proposal
		a.AcceptedProposal = proposal
		a.AcceptedValue = value
		// Accept proposal
		fmt.Printf("    node %d: STATUS   - %#v\n", a.Id, a)
		return AcceptResponse{
			Id:       a.Id,
			OK:       true,
			Proposal: proposal,
		}
	}
	fmt.Printf("    node %d: REJECTED - proposal:%d < PromisedProposal:%d (should be >=)\n", a.Id, proposal, a.PromisedProposal)
	return AcceptResponse{
		Id: a.Id,
		OK: false,
	}
}
