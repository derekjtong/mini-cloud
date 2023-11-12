package main

import (
	"fmt"

	"github.com/derekjtong/paxos/node"
	"github.com/derekjtong/paxos/utils"
)

func main() {
	// Initialize nodes based on configuration specified in utils/config.go
	for i, config := range utils.NodeConfigs {
		// Start Goroutine
		go func(config utils.NodeConfig) {
			fmt.Printf("Starting Node %d on %s:%d\n", i+1, config.IPAddress, config.Port)
			node := node.NewNode(config.IPAddress, config.Port)
			node.Start()
		}(config)
	}

	// Keep main function running to keep Goroutines running
	select {}
}
