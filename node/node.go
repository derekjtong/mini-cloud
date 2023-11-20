// node/node.go

package node

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

type Node struct {
	addr          string
	NodeID        int
	rpcClients    map[string]*rpc.Client
	NeighborNodes []string
}

func NewNode(nodeID int, addr string) *Node {
	return &Node{
		NodeID: nodeID,
		addr:   addr,
	}
}

// Ping
type PingRequest struct{}
type PingResponse struct {
	Message string
}

func (n *Node) Ping(req *PingRequest, res *PingResponse) error {
	fmt.Printf("[Node %d]: Pinged\n", n.NodeID)
	res.Message = "Pong from node " + strconv.Itoa(n.NodeID)
	return nil
}

// SetNeighbors
type SetNeighborsRequest struct {
	Neighbors []string
}
type SetNeighborsResponse struct {
}

func (n *Node) SetNeighbors(req *SetNeighborsRequest, res *SetNeighborsResponse) error {
	n.NeighborNodes = req.Neighbors
	n.rpcClients = make(map[string]*rpc.Client)

	for _, neighbor := range req.Neighbors {
		client, err := rpc.Dial("tcp", neighbor)
		if err != nil {
			fmt.Printf("[Node %d]: Error connecting to neighbor at %s: %v\n", n.NodeID, neighbor, err)
			continue
		}
		n.rpcClients[neighbor] = client
	}
	fmt.Printf("[Node %d]: Neighbors have been set successfully.\n", n.NodeID)
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

func (n *Node) Start() {
	listener, err := net.Listen("tcp", n.addr)
	if err != nil {
		fmt.Printf("[Node %d]: Error starting RPC server on %s: %v\n", n.NodeID, n.addr, err)
		return
	}

	defer listener.Close()

	fmt.Printf("[Node %d]: RPC server started on %s\n", n.NodeID, n.addr)

	rpcServer := rpc.NewServer()
	err = rpcServer.Register(n)
	if err != nil {
		fmt.Printf("[Node %d]: Error registering RPC server: %v\n", n.NodeID, err)
		return
	}

	rpcServer.Accept(listener)
	fmt.Printf("[Node %d]: Started on %s\n", n.NodeID, n.addr)
}
