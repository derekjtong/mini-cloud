// node/node.go

package node

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"

	"github.com/derekjtong/paxos/paxos"
)

type Node struct {
	addr          string
	NodeID        int
	rpcClients    map[string]*rpc.Client
	NeighborNodes []string
	proposer      *paxos.Proposer
	acceptor      *paxos.Acceptor
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
	fmt.Printf("[Node %d]: Creating directory %s\n", n.NodeID, fsDir)

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

	fmt.Printf("[Node %d]: Starting RPC server on %s\n", n.NodeID, n.addr)
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
	fmt.Printf("--------------------\n")
	fmt.Printf("[Node %d]: Client trying to write %s, running Paxos...\n", n.NodeID, req.Body)

	err := n.proposer.Propose(req.Body)
	if err != nil {
		return fmt.Errorf("could not achieve consensus")
	}

	// REDUNDANT: Done by local instance of Acceptor
	// Write Locally
	// if err := n.writeFileToLocal(n.proposer.Value); err != nil {
	// 	return fmt.Errorf("error writing file locally %v", err)
	// }
	// fmt.Printf("[Node %d]: PROPOSER - Wrote %s (proposer value)\n", n.NodeID, n.proposer.Value)
	fmt.Printf("--------------------\n")
	return nil
}

func (n *Node) writeFileToLocal(data string) error {
	filePath := fmt.Sprintf("./node_data/node_data_%s.json", n.addr)
	// file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
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
	filePath := fmt.Sprintf("./node_data/node_data_%s/data.json", n.addr)
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

	fmt.Printf("[Node %d]: WARNING, NO PAXOS. Read %s\n", n.NodeID, data)
	res.Data = data
	return nil
}

// RPC: Prepare
func (n *Node) Prepare(req *paxos.PrepareRequest, res *paxos.PrepareResponse) error {
	fmt.Printf("    node %d: RECEIVED - %#v\n", n.NodeID, req)
	// fmt.Printf("    node %d: received prepare request from %d {proposal: %d}\n", n.NodeID, req.Id, req.Proposal)
	*res = n.acceptor.Prepare(req.Proposal)
	// fmt.Printf("[Node %d]: Completed prepare\n", n.NodeID)
	return nil
}

// RPC: Accept
func (n *Node) Accept(req *paxos.AcceptRequest, res *paxos.AcceptResponse) error {
	fmt.Printf("    node %d: RECEIVED - %#v\n", n.NodeID, req)

	// fmt.Printf("    node %d: received accept request from %d {Proposal: %d, Value=%s}\n", n.NodeID, req.Id, req.Proposal, req.Value)

	*res = n.acceptor.Accept(req.Proposal, req.Value)
	n.acceptor.AcceptedValue = req.Value
	if res.OK {
		fmt.Printf("[Node %d]: ACCEPTOR - Wrote %s\n", n.NodeID, req.Value)
		n.writeFileToLocal(req.Value)
	}
	return nil
}

// RPC: info
type InfoRequest struct{}
type InfoResponse struct {
	ProposerInfo string
	AcceptorInfo string
}

func (n *Node) Info(req *InfoRequest, res *InfoResponse) error {
	// res.AcceptorInfo = fmt.Sprintf("%#v\n", n.acceptor)
	// res.ProposerInfo = fmt.Sprintf("%#v\n", n.proposer)

	res.AcceptorInfo = fmt.Sprintf("Acceptor={PromisedProposal:%d, AcceptedProposal:%d, AcceptedValue:%s}", n.acceptor.PromisedProposal, n.acceptor.AcceptedProposal, n.acceptor.AcceptedValue)
	res.ProposerInfo = fmt.Sprintf("Proposer={ProposalNumber:%d, Value:%s, HighestAcceptedProposalNumber:%d, HighestAcceptedValue:%s}", n.proposer.ProposalNumber, n.proposer.Value, n.proposer.HighestAcceptedProposalNumber, n.proposer.HighestAcceptedValue)
	return nil
}
