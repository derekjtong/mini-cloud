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
	fmt.Printf("    node %d: status   - %#v\n", a.Id, a)
	if proposal > a.PromisedProposal {
		fmt.Printf("%s    node %d: ACCEPTED - changing PromisedProposal from %d to incoming %d%s\n", Green, a.Id, a.PromisedProposal, proposal, Reset)
		a.PromisedProposal = proposal
		fmt.Printf("    node %d: status   - %#v\n", a.Id, a)
		// fmt.Printf("    node %d: {promisedProposal=%d, acceptedValue=%s, acceptedProposal=%d}\n", a.Id, a.PromisedProposal, a.AcceptedValue, a.AcceptedProposal)
		// Promise to not accept any earlier proposals
		return PrepareResponse{
			Id:            a.Id,
			OK:            true,
			Proposal:      a.AcceptedProposal,
			AcceptedValue: a.AcceptedValue,
		}
	}
	fmt.Printf("%s    node %d: REJECTED - PromisedProposal:%d greater than or equal to incoming proposal:%d%s\n", Red, a.Id, a.PromisedProposal, proposal, Reset)
	return PrepareResponse{
		Id: a.Id,
		OK: false,
	}
}

// Handle Accept request
func (a *Acceptor) Accept(proposal int, value string) AcceptResponse {
	fmt.Printf("    node %d: status   - %#v\n", a.Id, a)

	if proposal >= a.PromisedProposal {
		fmt.Printf("%s    node %d: ACCEPTED - proposal:%d >= PromisedProposal:%d%s\n", Green, a.Id, proposal, a.PromisedProposal, Reset)
		fmt.Printf("    node %d: updating - value '%s' -> '%s', AcceptedProposal %d -> %d\n", a.Id, a.AcceptedValue, value, a.AcceptedProposal, proposal)
		a.PromisedProposal = proposal
		a.AcceptedProposal = proposal
		a.AcceptedValue = value
		// Accept proposal
		fmt.Printf("    node %d: status   - %#v\n", a.Id, a)
		return AcceptResponse{
			Id:       a.Id,
			OK:       true,
			Proposal: proposal,
		}
	}
	fmt.Printf("%s    node %d: REJECTED - proposal:%d < PromisedProposal:%d (should be >=)%s\n", Red, a.Id, proposal, a.PromisedProposal, Reset)
	return AcceptResponse{
		Id: a.Id,
		OK: false,
	}
}
