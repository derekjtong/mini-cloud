// node/node.go

package node

import (
	"fmt"
	"net/rpc"

	myRPC "github.com/derekjtong/paxos/rpc"
)

type Node struct {
	Address       string
	NodeID        int
	rpcServer     *myRPC.RPCServer
	NeighborNodes []string
}

func NewNode(nodeID int, addr string) *Node {
	return &Node{
		NodeID:  nodeID,
		Address: addr,
	}
}

func (n *Node) AddNeighborNode(addr string) error {
	n.NeighborNodes = append(n.NeighborNodes, addr)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("error dialing %s: %v", addr, err)
	}
	defer client.Close()
	var request myRPC.PingRequest
	var response myRPC.PingResponse
	if err := client.Call("RPCServer.Ping", &request, &response); err != nil {
		return fmt.Errorf("error calling RPC method on %s: %v", addr, err)
	}

	fmt.Printf("%+v. Connected!\n", response.Message)
	return nil
}

func (n *Node) Start() {
	// Start the RPC server
	n.rpcServer = myRPC.NewServer(n.NodeID, n.Address)
	go n.rpcServer.Start()

	// Now, you can perform any additional initialization or start other components
	// ...

	// Print the final message
	fmt.Printf("[Node %d]: Started on %s\n", n.NodeID, n.Address)
}
