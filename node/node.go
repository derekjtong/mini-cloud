// node/node.go

package node

import (
	"fmt"

	myRPC "github.com/derekjtong/paxos/rpc"
)

type Node struct {
	Address       string
	NodeID        int
	rpcServer     *myRPC.RPCServer
}

func NewNode(nodeID int, addr string) *Node {
	return &Node{
		NodeID:  nodeID,
		Address: addr,
	}
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
