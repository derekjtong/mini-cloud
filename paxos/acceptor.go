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
	fmt.Printf("    node %d: BEFORE %#v\n", a.Id, a)
	if proposal > a.PromisedProposal {
		fmt.Printf("    node %d: accepted prepare preposal, (current)promisedProposal=%d, (incoming)proposal=%d\n", a.Id, a.PromisedProposal, proposal)
		a.PromisedProposal = proposal
		fmt.Printf("    node %d: AFTER %#v\n", a.Id, a)
		// fmt.Printf("    node %d: {promisedProposal=%d, acceptedValue=%s, acceptedProposal=%d}\n", a.Id, a.PromisedProposal, a.AcceptedValue, a.AcceptedProposal)
		// Promise to not accept any earlier proposals
		return PrepareResponse{
			Id:            a.Id,
			OK:            true,
			Proposal:      a.AcceptedProposal,
			AcceptedValue: a.AcceptedValue,
		}
	}
	fmt.Printf("    node %d: rejected prepare proposal {promisedProposal=%d}\n", a.Id, a.PromisedProposal)
	return PrepareResponse{
		Id: a.Id,
		OK: false,
	}
}

// Handle Accept request
func (a *Acceptor) Accept(proposal int, value string) AcceptResponse {
	if proposal >= a.PromisedProposal {
		a.PromisedProposal = proposal
		a.AcceptedProposal = proposal
		a.AcceptedValue = value
		fmt.Printf("    node %d: promised proposal=%d\n", a.Id, a.PromisedProposal)
		fmt.Printf("    node %d: accepted proposal=%d\n", a.Id, a.AcceptedProposal)
		fmt.Printf("    node %d: accepted value %s\n", a.Id, a.AcceptedValue)
		// Accept proposal
		return AcceptResponse{
			Id:       a.Id,
			OK:       true,
			Proposal: proposal,
		}
	}
	fmt.Printf("Rejected")
	return AcceptResponse{
		Id: a.Id,
		OK: false,
	}
}
