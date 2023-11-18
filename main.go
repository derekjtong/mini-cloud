package main

import (
	"fmt"

	"github.com/derekjtong/paxos/node"
	"github.com/derekjtong/paxos/utils"
)

func main() {
	for i, config := range utils.NodeConfigs {
		// Start Goroutine
		go func(i int, config utils.NodeConfig) {
			fmt.Printf("Starting Node %d on %s:%d\n", i+1, config.IPAddress, config.Port)
			node := node.NewNode(config.IPAddress, config.Port)
			node.Start()
		}(i, config)
	}

	// Keep main function running to keep Goroutines running
	select {}
}
