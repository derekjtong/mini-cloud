// node/node.go

package node

import (
	"fmt"
	"net/rpc"
	"os"

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

func (n *Node) AddNeighborNode(addr string) {
	n.NeighborNodes = append(n.NeighborNodes, addr)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Error dialing!")
	}
	defer client.Close()
	var request myRPC.PingRequest
	var response myRPC.PingResponse
	if err := client.Call("RPCServer.Ping", &request, &response); err != nil {
		fmt.Printf("Error calling RPC method: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v. Connected!\n", response.Message)
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
