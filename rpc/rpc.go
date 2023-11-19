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
	addr   string
	NodeID int
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

func NewServer(nodeID int, addr string) *RPCServer {
	return &RPCServer{
		NodeID: nodeID,
		addr:   addr,
	}
}

func (s *RPCServer) Start() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Printf("[Node %d]: Error starting RPC server on %s: %v\n", s.NodeID, s.addr, err)
		return
	}

	defer listener.Close()

	fmt.Printf("[Node %d]: RPC server started on %s\n", s.NodeID, s.addr)

	rpcServer := rpc.NewServer()
	err = rpcServer.Register(s)
	if err != nil {
		fmt.Printf("[Node %d]: Error registering RPC server: %v\n", s.NodeID, err)
		return
	}

	rpcServer.Accept(listener)
}
