// main.go

package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"

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

	// fmt.Printf("%v, connected!", response.Message)
	fmt.Printf("Connected to node %v\n", response.NodeID)

	runCLI(client)
}
func startServer() {
	fmt.Printf("Starting server! Hint: to start client, 'go run main.go client'.\n\n")

	var nodeAddrList []string

	// Start 3 nodes
	for nodeID := 1; nodeID <= utils.NodeCount; nodeID++ {
		port, err := findAvailablePort()
		if err != nil {
			fmt.Printf("Error finding available port: %v\n", err)
			return
		}
		addr := fmt.Sprintf("%s:%d", utils.IPAddress, port)
		nodeAddrList = append(nodeAddrList, addr)
		go func(addr string, nodeNumber int) {
			fmt.Printf("[Node %d]: Starting on %s\n", nodeNumber, addr)
			node, err := node.NewNode(nodeNumber, addr)
			if err != nil {
				fmt.Printf("Error creating node %d: %v", nodeID, err)
				return
			}
			node.Start()
		}(addr, nodeID)
		// Wait until server is ready
		err = waitForServerReady(addr)
		if err != nil {
			fmt.Printf("Error waiting for node %d to be ready: %v\n", nodeID, err)
			return
		}
	}
	// Send list of IP addresses to nodes
	for _, nodeAddr := range nodeAddrList {
		client, err := rpc.Dial("tcp", nodeAddr)
		if err != nil {
			fmt.Printf("[SERVER] Error dialing node %s: %v\n", nodeAddr, err)
			continue
		}
		var setNeighborsRequest = node.SetNeighborsRequest{Neighbors: nodeAddrList}
		var setNeighborsResponse node.SetNeighborsResponse
		if err := client.Call("Node.SetNeighbors", &setNeighborsRequest, &setNeighborsResponse); err != nil {
			fmt.Printf("Error setting neighbors for node %s: %v\n", nodeAddr, err)
		}
		client.Close()
	}
	select {}
}

// Check server is ready
func waitForServerReady(address string) error {
	// Exponential backoff
	var backoff time.Duration = 100
	const maxBackoff = 5 * time.Second
	const maxRetries = 10
	var timeout time.Duration = 5 * time.Second
	startTime := time.Now()

	for retries := 0; retries < maxRetries; retries++ {
		client, err := rpc.Dial("tcp", address)
		if err == nil {
			var req node.HealthCheckRequest
			var res node.HealthCheckResponse
			err = client.Call("Node.HealthCheck", &req, &res)
			client.Close()
			if err == nil && res.Status == "OK" {
				return nil
			}
		}

		if time.Since(startTime) > timeout {
			return fmt.Errorf("server at %s did not become ready within %v", address, timeout)
		}

		if backoff < maxBackoff {
			backoff *= 2
		}
		time.Sleep(backoff)
	}
	return fmt.Errorf("server at %s did not become ready afte %d attemps", address, maxRetries)
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

func runCLI(client *rpc.Client) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commands (type 'exit' to quit):")

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		parts := strings.SplitN(input, " ", 2)
		command := parts[0]
		var argument string
		if len(parts) > 1 {
			argument = parts[1]
		}

		// Process commands
		switch command {
		case "ping":
			var req node.PingRequest
			var res node.PingResponse
			if err := client.Call("Node.Ping", &req, &res); err != nil {
				fmt.Printf("Error calling RPC method: %v\n", err)
				continue
			}
			fmt.Println(res.Message)
		case "write":
			if argument == "" {
				fmt.Println("Please provide a string to write")
				continue
			}
			var req node.WriteFileRequest
			var res node.WriteFileResponse
			req.Body = argument
			if err := client.Call("Node.WriteFile", &req, &res); err != nil {
				fmt.Printf("Error calling RPC method: %v\n", err)
			} else {
				fmt.Println("Write operation successful")
			}
		default:
			fmt.Println("Unknown command:", input)
		}
	}
}
