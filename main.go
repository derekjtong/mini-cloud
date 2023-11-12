package main

import (
	"fmt"

	"github.com/derekjtong/paxos/utils"
)

func main() {
	fmt.Println("Hello")
	node1IP := utils.NodeConfigs[0].IPAddress
	node1Port := utils.NodeConfigs[0].Port
	fmt.Printf("%s:%d\n", node1IP, node1Port)
}
