// rpc/server.go

package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// RPCServer implementation

type RPCServer struct {
	IPAddress string
	Port      int
	NodeID    int
	wg        *sync.WaitGroup // Added a WaitGroup
}

func NewServer(nodeID int, ipAddress string, port int, wg *sync.WaitGroup) *RPCServer {
	return &RPCServer{
		NodeID:    nodeID,
		IPAddress: ipAddress,
		Port:      port,
		wg:        wg,
	}
}

func (s *RPCServer) Start() {
	addr := fmt.Sprintf("%s:%d", s.IPAddress, s.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("[Node%d]: Error starting RPC server on %s: %v\n", s.NodeID, addr, err)
		return
	}

	defer listener.Close()

	fmt.Printf("[Node%d]: RPC server started on %s\n", s.NodeID, addr)

	// Signal that the RPC server has started
	s.wg.Done()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[Node%d]: Error accepting connection: %v\n", s.NodeID, err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
