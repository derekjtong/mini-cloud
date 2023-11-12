package node

import (
	"fmt"

	"github.com/derekjtong/paxos/rpc"
)

type Node struct {
	IPAddress string
	Port      int
}

func NewNode(ipAddress string, port int) *Node {
	return &Node{
		IPAddress: ipAddress,
		Port:      port,
	}
}

func (n *Node) Start() {
	rpcServer := rpc.NewServer(n.IPAddress, n.Port)
	go rpcServer.Start()
	fmt.Printf("Node started on %s:%d\n", n.IPAddress, n.Port)
}
