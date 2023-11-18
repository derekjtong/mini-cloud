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

	for i, config := range utils.NodeConfigs {
		// Dynamically find an available port
		port, err := findAvailablePort()
		if err != nil {
			fmt.Printf("Error finding available port: %v\n", err)
			return
		}

		// Start Goroutine
		wg.Add(1)
		go func(i int, config utils.NodeConfig, port int) {
			defer wg.Done()

			fmt.Printf("Starting Node %d on %s:%d\n", i+1, config.IPAddress, port)
			node := node.NewNode(config.IPAddress, port)
			node.Start()
		}(i, config, port)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
}
