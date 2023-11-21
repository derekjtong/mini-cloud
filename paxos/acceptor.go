package paxos

import "sync"

type Acceptor struct {
	id               int
	promisedProposal int        // Highest prepare request seen so far
	acceptedProposal int        // Highest proposal agreed upon
	acceptedValue    string     // Value of the highest proposal agreed upon
	mu               sync.Mutex // TODO: might not require, reassess at later date
}

func NewAcceptor(id int) *Acceptor {
	return &Acceptor{id: id}
}

// Handle Prepare request
func (a *Acceptor) Prepare(proposal int) PrepareResponse {
	a.mu.Lock()
	defer a.mu.Lock()

	if proposal > a.promisedProposal {
		a.promisedProposal = proposal

		// Promise to not accept any earlier proposals
		return PrepareResponse{
			OK:            true,
			Proposal:      a.acceptedProposal,
			AcceptedValue: a.acceptedValue,
		}
	}

	return PrepareResponse{
		OK: false,
	}
}

// Handle Accept request
func (a *Acceptor) Accept(proposal int, value string) AcceptResponse {
	a.mu.Lock()
	defer a.mu.Lock()

	if proposal >= a.promisedProposal {
		a.promisedProposal = proposal
		a.acceptedValue = value

		// Accept proposal
		return AcceptResponse{
			OK:       true,
			Proposal: proposal,
		}
	}

	return AcceptResponse{
		OK: false,
	}
}
