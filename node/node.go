// node/node.go

package node

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/derekjtong/mini-cloud/paxos"
	"github.com/derekjtong/mini-cloud/utils"
)

type Node struct {
	addr          string
	NodeID        int
	rpcClients    map[string]*rpc.Client
	NeighborNodes []string
	proposer      *paxos.Proposer
	acceptor      *paxos.Acceptor
	terminated    bool
	stop          bool
}

func NewNode(nodeID int, addr string) (*Node, error) {
	if addr == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	acceptor := paxos.NewAcceptor(nodeID)
	return &Node{
		NodeID:        nodeID,
		addr:          addr,
		rpcClients:    make(map[string]*rpc.Client),
		NeighborNodes: make([]string, 0),
		acceptor:      acceptor,
		stop:          false,
		// proposer initialized under SetNeighbors
	}, nil
}

// Start
func (n *Node) Start() {
	fsDir := "./node_data/"
	if err := os.MkdirAll(fsDir, 0755); err != nil {
		fmt.Printf("[Node %d]: Error creating file system directory: %v\n", n.NodeID, err)
		return
	}
	if !utils.MinimalStartUpLogging {
		fmt.Printf("[Node %d]: Creating directory %s\n", n.NodeID, fsDir)
	}

	listener, err := net.Listen("tcp", n.addr)
	if err != nil {
		fmt.Printf("[Node %d]: Error starting RPC server on %s: %v\n", n.NodeID, n.addr, err)
		return
	}

	defer listener.Close()

	rpcServer := rpc.NewServer()
	err = rpcServer.Register(n)
	if err != nil {
		fmt.Printf("[Node %d]: Error registering RPC server: %v\n", n.NodeID, err)
		return
	}
	if !utils.MinimalStartUpLogging {
		fmt.Printf("[Node %d]: Starting RPC server on %s\n", n.NodeID, n.addr)
	}
	rpcServer.Accept(listener)
}

// RPC: Ping
type PingRequest struct{}
type PingResponse struct {
	Message string
	NodeID  int
}

func (n *Node) Ping(req *PingRequest, res *PingResponse) error {
	fmt.Printf("[Node %d]: Pinged\n", n.NodeID)
	res.Message = "Pong from node " + strconv.Itoa(n.NodeID)
	res.NodeID = n.NodeID
	return nil
}

// PRC: SetNeighbors
type SetNeighborsRequest struct {
	Neighbors []string
}
type SetNeighborsResponse struct {
}

// Update node's list of neighbors
func (n *Node) SetNeighbors(req *SetNeighborsRequest, res *SetNeighborsResponse) error {
	n.NeighborNodes = req.Neighbors
	for _, neighbor := range req.Neighbors {
		// Check to not include node's own IP address
		// if neighbor != n.addr && n.rpcClients[neighbor] == nil {
		client, err := rpc.Dial("tcp", neighbor)
		if err != nil {
			fmt.Printf("[Node %d]: Error connecting to neighbor at %s: %v\n", n.NodeID, neighbor, err)
			continue
		}
		n.rpcClients[neighbor] = client
		// }
	}
	fmt.Printf("[Node %d]: Set neighbors\n", n.NodeID)

	// Initialize proposer
	n.proposer = paxos.NewProposer(n.NodeID, n.NodeID, n.rpcClients)

	return nil
}

// RPC: Health Check
type HealthCheckRequest struct{}
type HealthCheckResponse struct {
	Status string
}

func (n *Node) HealthCheck(req *HealthCheckRequest, res *HealthCheckResponse) error {
	res.Status = "OK"
	return nil
}

// RPC: WriteFile
type WriteFileRequest struct {
	Body string
}
type WriteFileResponse struct {
}

func (n *Node) WriteFile(req *WriteFileRequest, res *WriteFileResponse) error {
	fmt.Printf("[Node %d]: Client trying to write %s, running Paxos...\n", n.NodeID, req.Body)
	fmt.Printf("--------------------\n")

	err := n.proposer.Propose(req.Body)
	if err != nil {
		return err
	}

	fmt.Printf("[Node %d]: Paxos completed\n", n.NodeID)
	fmt.Printf("--------------------\n")
	return nil
}

// RPC: ForceWrite - WriteFile with retry
func (n *Node) ForceWrite(req *WriteFileRequest, res *WriteFileResponse) error {
	fmt.Printf("[Node %d]: Client trying to write %s, running Paxos...\n", n.NodeID, req.Body)
	fmt.Printf("--------------------\n")

	const maxRetries = 5
	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			fmt.Printf("[Node %d]: Retrying attempt %d/%d\n", n.NodeID, attempt, maxRetries)
			// Randomized delay
			r := rand.Intn(5-1+1) + 1
			time.Sleep(time.Duration(r) * time.Second)
		}

		err = n.proposer.Propose(req.Body)
		if err == nil {
			fmt.Printf("[Node %d]: Paxos completed successfully\n", n.NodeID)
			fmt.Printf("--------------------\n")
			return nil
		}
	}

	fmt.Printf("[Node %d]: Paxos failed after %d attempts: %v\n", n.NodeID, maxRetries, err)
	fmt.Printf("--------------------\n")
	return fmt.Errorf("could not achieve consensus after %d attempts: %v", maxRetries, err)
}

func (n *Node) writeFileToLocal(data string) error {
	filePath := fmt.Sprintf("./node_data/node_data_%s.json", n.addr)
	// file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) //Append
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) // Overwrite
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return nil
}

// RPC: ReadFile
type ReadFileRequest struct {
}

type ReadFileResponse struct {
	Data string
}

func (n *Node) ReadFile(req *ReadFileRequest, res *ReadFileResponse) error {
	filePath := fmt.Sprintf("./node_data/node_data_%s.json", n.addr)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data string
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	fmt.Printf("[Node %d]: Read %s\n", n.NodeID, data)
	res.Data = data
	return nil
}

// RPC: Prepare
func (n *Node) Prepare(req *paxos.PrepareRequest, res *paxos.PrepareResponse) error {
	if n.stop {
		return nil
	}
	fmt.Printf("    node %d: received - %#v\n", n.NodeID, req)
	*res = n.acceptor.Prepare(req.Proposal)
	return nil
}

// RPC: Accept
func (n *Node) Accept(req *paxos.AcceptRequest, res *paxos.AcceptResponse) error {
	if n.stop {
		return nil
	}
	fmt.Printf("    node %d: received - %#v\n", n.NodeID, req)

	*res = n.acceptor.Accept(req.Proposal, req.Value)
	n.acceptor.AcceptedValue = req.Value
	if res.OK {
		fmt.Printf("[Node %d]: ACCEPTOR - Wrote %s\n", n.NodeID, req.Value)
		n.writeFileToLocal(req.Value)
	}
	return nil
}

// RPC: Info
type InfoRequest struct{}
type InfoResponse struct {
	ProposerInfo string
	AcceptorInfo string
}

func (n *Node) Info(req *InfoRequest, res *InfoResponse) error {
	// Not used because too verbose: will print out *rpc.Clients map
	// res.AcceptorInfo = fmt.Sprintf("%#v\n", n.acceptor)
	// res.ProposerInfo = fmt.Sprintf("%#v\n", n.proposer)

	res.AcceptorInfo = fmt.Sprintf("Acceptor={PromisedProposal:%d, AcceptedProposal:%d, AcceptedValue:\"%s\"}", n.acceptor.PromisedProposal, n.acceptor.AcceptedProposal, n.acceptor.AcceptedValue)
	res.ProposerInfo = fmt.Sprintf("Proposer={ProposalNumber:%d, Value:\"%s\", HighestAcceptedProposalNumber:%d, HighestAcceptedValue:%s}", n.proposer.ProposalNumber, n.proposer.Value, n.proposer.HighestAcceptedProposalNumber, n.proposer.HighestAcceptedValue)
	return nil
}

// RPC: Toggletimeout
type TimeoutRequest struct{}
type TimeoutResponse struct {
	IsTimeout bool
}

func (n *Node) ToggleTimeout(req *TimeoutRequest, res *TimeoutResponse) error {
	if n.proposer.Timeout {
		n.proposer.Timeout = false
	} else {
		n.proposer.Timeout = true
	}

	res.IsTimeout = n.proposer.Timeout
	fmt.Printf("[Node %d]: Timeout ", n.NodeID)
	if n.proposer.Timeout {
		fmt.Printf("on\n")
	} else {
		fmt.Printf("off\n")
	}
	return nil
}

// RPC: Ping
type StopRequest struct{}
type StopResponse struct {
	IsStopped bool
}

func (n *Node) ToggleStop(req *StopRequest, res *StopResponse) error {
	n.stop = !n.stop
	fmt.Printf("[Node %d]: Client toggled stop, ", n.NodeID)
	if n.stop {
		fmt.Printf("server will no longer respond to Paxos\n")
	} else {
		fmt.Printf("server will respond to Paxos\n")
	}
	res.IsStopped = n.stop
	return nil
}

// RPC: Terminate
type TerminateRequest struct{}
type TerminateResponse struct{}

func (n *Node) Terminate(req *TerminateRequest, res *TerminateResponse) error {
	fmt.Printf("[Node %d]: Terminate method called\n", n.NodeID)

	// Avoid repeated termination
	if n.terminated {
		return nil
	}

	// Set termination flag
	n.terminated = true

	// Only send Terminate RPC to neighbors
	for neighborAddr, client := range n.rpcClients {
		if neighborAddr != n.addr {
			var terminateRequest TerminateRequest
			var terminateResponse TerminateResponse
			if err := client.Call("Node.Terminate", &terminateRequest, &terminateResponse); err != nil {
				fmt.Printf("[Node %d]: Error calling Terminate RPC method to %s: %v\n", n.NodeID, neighborAddr, err)
			}
		}
	}

	os.Exit(0)
	return nil
}
