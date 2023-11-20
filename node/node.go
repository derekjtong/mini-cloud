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

func (n *Node) SetNeighborNodes(neighbors []string) error {
	return nil
}

type PingRequest struct{}
type PingResponse struct {
	Message string
}

func (n *Node) Ping(req *PingRequest, res *PingResponse) error {
	fmt.Printf("[Node %d]: Pinged\n", n.NodeID)
	res.Message = "Pong from node " + strconv.Itoa(n.NodeID)
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
