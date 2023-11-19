// rpc/rpc.go

package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

// RPCServer implementation

type RPCServer struct {
	IPAddress string
	Port      int
	NodeID    int
}

type PingRequest struct{}
type PingResponse struct {
	Message string
}

func (s *RPCServer) Ping(req *PingRequest, res *PingResponse) error {
	fmt.Printf("[Node %d]: Pinged\n", s.NodeID)
	res.Message = "Pong from node " + strconv.Itoa(s.NodeID)
	return nil
}

func NewServer(nodeID int, ipAddress string, port int) *RPCServer {
	return &RPCServer{
		NodeID:    nodeID,
		IPAddress: ipAddress,
		Port:      port,
	}
}

func (s *RPCServer) Start() {
	addr := fmt.Sprintf("%s:%d", s.IPAddress, s.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("[Node %d]: Error starting RPC server on %s: %v\n", s.NodeID, addr, err)
		return
	}

	defer listener.Close()

	fmt.Printf("[Node %d]: RPC server started on %s\n", s.NodeID, addr)

	rpcServer := rpc.NewServer()
	err = rpcServer.Register(s)
	if err != nil {
		fmt.Printf("[Node %d]: Error registering RPC server: %v\n", s.NodeID, err)
		return
	}

	rpcServer.Accept(listener)
}
