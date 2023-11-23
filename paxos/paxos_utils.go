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

// Colors for terminal
const (
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Reset  = "\033[0m"
)
