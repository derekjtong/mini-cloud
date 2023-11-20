// node/node.go

package node

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

type Node struct {
	addr          string
	NodeID        int
	rpcClients    map[string]*rpc.Client
	NeighborNodes []string
}

func NewNode(nodeID int, addr string) (*Node, error) {
	if addr == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}
	return &Node{
		NodeID:        nodeID,
		addr:          addr,
		rpcClients:    make(map[string]*rpc.Client),
		NeighborNodes: make([]string, 0),
	}, nil
}

// Start
func (n *Node) Start() {
	fsDir := fmt.Sprintf("./node_data/node_data_%s", n.addr)
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

// Ping
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

// SetNeighbors
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
		if neighbor != n.addr && n.rpcClients[neighbor] == nil {
			client, err := rpc.Dial("tcp", neighbor)
			if err != nil {
				fmt.Printf("[Node %d]: Error connecting to neighbor at %s: %v\n", n.NodeID, neighbor, err)
				continue
			}
			n.rpcClients[neighbor] = client
		}
	}
	fmt.Printf("[Node %d]: Set neighbors\n", n.NodeID)
	return nil
}

// Health Check
type HealthCheckRequest struct{}
type HealthCheckResponse struct {
	Status string
}

func (n *Node) HealthCheck(req *HealthCheckRequest, res *HealthCheckResponse) error {
	res.Status = "OK"
	return nil
}

// WriteFile
type WriteFileRequest struct {
	Body string
}
type WriteFileResponse struct {
}

func (n *Node) WriteFile(req *WriteFileRequest, res *WriteFileResponse) error {
	// RunPaxos()

	filePath := fmt.Sprintf("./node_data/node_data_%s/data.json", n.addr)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(req.Body); err != nil {
		return err
	}
	return nil
}
