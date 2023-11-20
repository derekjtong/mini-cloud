// main.go

package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"

	"github.com/derekjtong/paxos/node"
	"github.com/derekjtong/paxos/utils"
)

type NodeIPs struct {
	IPs []string
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "client" {
		startClient()
	} else {
		startServer()
	}

}

func startClient() {
	fmt.Print("Starting Client!\nNode IP address: (defaulting to 127.0.0.1)\n")
	var IPAddress string = "127.0.0.1"
	// fmt.Scanln(&IPAddress)
	fmt.Print("Node port number: ")
	var Port int
	fmt.Scanln(&Port)
	fmt.Printf("Connecting to %s:%d...\n", IPAddress, Port)
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", IPAddress, Port))
	if err != nil {
		fmt.Printf("Error dialing RPC server:%v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	var request node.PingRequest
	var response node.PingResponse
	if err := client.Call("Node.Ping", &request, &response); err != nil {
		fmt.Printf("Error calling RPC method: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Reply: %+v\n", response.Message)
}
func startServer() {
	fmt.Printf("Starting server! Hint: to start client, 'go run main.go client'.\n\n")
	var wg sync.WaitGroup

	for nodeID := 1; nodeID <= utils.NodeCount; nodeID++ {
		wg.Add(1)
		port, err := findAvailablePort()
		if err != nil {
			fmt.Printf("Error finding available port: %v\n", err)
			return
		}
		addr := fmt.Sprintf("%s:%d", utils.IPAddress, port)
		go func(addr string, nodeNumber int, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Printf("[Node %d]: Starting on %s\n", nodeNumber, addr)
			node := node.NewNode(nodeNumber, addr)
			node.Start()
		}(addr, nodeID, &wg)
	}
	// Send NodeNeighbors to every node
	// for _, addr := range nodeNeighbors {
	// client, err := rpc.Dial("tcp", addr)
	// if err != nil {
	// 	fmt.Printf("Error dialing node %s: %v\n", addr, err)
	// 	continue
	// }
	// var setNeighborsRequest = myRPC.SetNeighborNodes{Neighbors: nodeNeighbors}
	// var setNeighborsResponse myRpc.setNeighborsResponse
	// if err := client.Call("Node.SetNeighborNodes", &setNeighborsRequest, &setNeighborsResponse); err != nil {
	// 	fmt.Printf("Error setting neighbors for node %s: %v\n", addr, err)
	// }
	// client.Close()
	// }

	wg.Wait()
	select {}
}

func findAvailablePort() (int, error) {
	// Find a free port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	// Get the allocated port
	address := listener.Addr().String()
	_, portString, err := net.SplitHostPort(address)
	if err != nil {
		return 0, err
	}

	port, err := net.LookupPort("tcp", portString)
	if err != nil {
		return 0, err
	}

	return port, nil
}
