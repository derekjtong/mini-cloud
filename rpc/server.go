package rpc

import (
	"fmt"
	"net"
	"net/rpc"
)

// RPC server implementation

type RPCServer struct {
	IPAddress string
	Port      int
}

func NewServer(ipAddress string, port int) *RPCServer {
	return &RPCServer{
		IPAddress: ipAddress,
		Port:      port,
	}
}

func (s *RPCServer) Start() {
	addr := fmt.Sprintf("%s:%d", s.IPAddress, s.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Error starting RPC server on %s: %v\n", addr, err)
		return
	}

	defer listener.Close()

	fmt.Printf("RPC server started on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
