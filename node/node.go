// node/node.go

package node

import (
	"fmt"
	"sync"

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

func (n *Node) Start(rpcWG *sync.WaitGroup) {
	// Start the RPC server
	n.rpcServer = rpc.NewServer(n.NodeID, n.IPAddress, n.Port, rpcWG)
	go n.rpcServer.Start()

	// Wait for the RPC server to start
	rpcWG.Wait()

	// Now, you can perform any additional initialization or start other components
	// ...

	// Print the final message
	fmt.Printf("[Node %d]: Started on %s:%d\n", n.NodeID, n.IPAddress, n.Port)
}
