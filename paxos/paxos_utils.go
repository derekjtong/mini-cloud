package paxos

// Prepare phase request
type PrepareRequest struct {
	Id       int
	Proposal int
}

// Prepare phase response
type PrepareResponse struct {
	Id            int
	OK            bool
	Proposal      int
	AcceptedValue string
}

// Accept phase request
type AcceptRequest struct {
	Id       int
	Proposal int
	Value    string
}

// Accept phase response
type AcceptResponse struct {
	Id       int
	OK       bool
	Proposal int
}
