// node/node.go

package node

import (
	"fmt"

	"github.com/derekjtong/paxos/rpc"
)

type Node struct {
	IPAddress string
	Port      int
	NodeID    int
	rpcServer *rpc.RPCServer // Added a field to store RPCServer instance
}

func NewNode(nodeID int, ipAddress string, port int) *Node {
	return &Node{
		NodeID:    nodeID,
		IPAddress: ipAddress,
		Port:      port,
	}
}

func (n *Node) Start() {
	// Start the RPC server
	n.rpcServer = rpc.NewServer(n.NodeID, n.IPAddress, n.Port)
	go n.rpcServer.Start()

	// Now, you can perform any additional initialization or start other components
	// ...

	// Print the final message
	fmt.Printf("[Node %d]: Started on %s:%d\n", n.NodeID, n.IPAddress, n.Port)
}
