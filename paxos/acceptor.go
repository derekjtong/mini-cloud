package paxos

import (
	"fmt"
	"sync"
)

type Acceptor struct {
	id               int
	promisedProposal int        // Highest prepare request seen so far
	acceptedProposal int        // Highest proposal agreed upon
	acceptedValue    string     // Value of the highest proposal agreed upon
	mu               sync.Mutex // TODO: might not require, reassess at later date
}

func NewAcceptor(id int) *Acceptor {
	return &Acceptor{
		id:               id,
		promisedProposal: 0,
		acceptedProposal: 0,
		acceptedValue:    "",
	}
}

// Handle Prepare request
func (a *Acceptor) Prepare(proposal int) PrepareResponse {
	a.mu.Lock()
	defer a.mu.Unlock()

	if proposal > a.promisedProposal {
		a.promisedProposal = proposal
		fmt.Printf("    node %d: accepted proposal\n", a.id)

		// Promise to not accept any earlier proposals
		return PrepareResponse{
			Id:            a.id,
			OK:            true,
			Proposal:      a.acceptedProposal,
			AcceptedValue: a.acceptedValue,
		}
	}
	fmt.Printf("    node %d: rejected proposal\n", a.id)
	return PrepareResponse{
		Id: a.id,
		OK: false,
	}
}

// Handle Accept request
func (a *Acceptor) Accept(proposal int, value string) AcceptResponse {
	a.mu.Lock()
	defer a.mu.Unlock()

	if proposal >= a.promisedProposal {
		fmt.Printf("    node %d: accepted\n", a.id)

		a.promisedProposal = proposal
		a.acceptedValue = value

		// Accept proposal
		return AcceptResponse{
			Id:       a.id,
			OK:       true,
			Proposal: proposal,
		}
	}

	return AcceptResponse{
		Id: a.id,
		OK: false,
	}
}
