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

func main() {
	if len(os.Args) > 1 && os.Args[1] == "client" {
		startClient()
	} else {
		startServer()
	}
}

func startClient() {
	fmt.Print("Starting Client!\n\nNode IP address: (defaulting to 127.0.0.1)\n")

	var ipAddress string = "127.0.0.1"
	// fmt.Scanln(&IPAddress)

	fmt.Print("Node port number: ")
	var port int
	fmt.Scanln(&port)

	fmt.Printf("Pinging %s:%d\n", ipAddress, port)
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", ipAddress, port))
	if err != nil {
		fmt.Printf("Error dialing RPC server, please confirm port number.\n            (%v)\n", err)
		os.Exit(1)
	}

	defer client.Close()

	var request node.PingRequest
	var response node.PingResponse
	if err := client.Call("Node.Ping", &request, &response); err != nil {
		fmt.Printf("Error calling RPC method: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v. Connected!\n", response.Message)
}

func startServer() {
	fmt.Printf("Starting server! Hint: to start client, 'go run main.go client'.\n\n")
	var wg sync.WaitGroup
	var nodeNeighbors []string
	for nodeID := 1; nodeID <= utils.NodeCount; nodeID++ {
		wg.Add(1)
		port, err := findAvailablePort()
		if err != nil {
			fmt.Printf("Error finding available port: %v\n", err)
			return
		}
		addr := fmt.Sprintf("%s:%d", utils.IPAddress, port)
		nodeNeighbors = append(nodeNeighbors, addr)
		go func(ipAddress string, nodeID int, port int, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Printf("[Node %d]: Starting on %s\n", nodeID, addr)
			node := node.NewNode(nodeID, addr)
			node.Start()
		}(utils.IPAddress, nodeID, port, &wg)
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

// Find an available port
func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

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
