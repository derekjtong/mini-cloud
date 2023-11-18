// main.go

package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/derekjtong/paxos/node"
	"github.com/derekjtong/paxos/utils"
)

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

func main() {
	var wg sync.WaitGroup
	var rpcWG sync.WaitGroup

	for _, config := range utils.NodeConfigs {
		// Increment the main wait group for each node
		wg.Add(1)

		// Increment the RPC wait group for each node
		rpcWG.Add(1)

		// Dynamically find an available port
		port, err := findAvailablePort()
		if err != nil {
			fmt.Printf("Error finding available port: %v\n", err)
			return
		}

		// Start Goroutine for node
		go func(config utils.NodeConfig, port int, wg *sync.WaitGroup, rpcWG *sync.WaitGroup) {
			defer wg.Done()

			fmt.Printf("[Node%d]: Starting on %s:%d\n", config.NodeID, config.IPAddress, port)
			node := node.NewNode(config.NodeID, config.IPAddress, port)

			// Start the node
			node.Start(rpcWG)

		}(config, port, &wg, &rpcWG)
	}

	// Wait for all nodes to finish starting
	wg.Wait()

	// Keep the main function running to keep the servers active
	select {}
}
