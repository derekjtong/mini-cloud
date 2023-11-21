package paxos

// Prepare phase request
type PrepareRequest struct {
	Proposal int
}

// Prepare phase response
type PrepareResponse struct {
	OK            bool
	Proposal      int
	AcceptedValue string
}

// Accept phase request
type AcceptRequest struct {
	Proposal int
	Value    string
}

// Accept phase response
type AcceptResponse struct {
	OK       bool
	Proposal int
}
