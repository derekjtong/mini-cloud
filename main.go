// main.go

package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/derekjtong/mini-cloud/node"
	"github.com/derekjtong/mini-cloud/utils"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "client":
			startClient()
		default:
			fmt.Printf("Invalid arg")
		}
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
		fmt.Printf("Error dialing RPC server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	var request node.PingRequest
	var response node.PingResponse
	if err := client.Call("Node.Ping", &request, &response); err != nil {
		fmt.Printf("Error calling RPC method: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connected to node %v!\n", response.NodeID)
	if os.Args[1] == "kill" {
		os.Setenv("TERMINATE", "true")
		fmt.Println("TERMINATE signal sent. Exiting.")
		os.Exit(0)
	}

	runCLI(client)
}

func startServer() {
	fmt.Printf("Starting server! Hint: to start client, 'go run main.go client'.\n\n")

	if utils.ClearNodeDataOnStart {
		fmt.Println("[SERVER]: Clearing node_data directory")
		// Create directory if it doesn't exist
		if err := os.MkdirAll("./node_data", 0755); err != nil {
			fmt.Printf("Error creating node_data directory: %v\n", err)
			return
		}
		if err := clearDir("./node_data"); err != nil {
			fmt.Printf("Error clearing node_data directory: %v\n", err)
			return
		}
	}

	var nodeAddrList []string

	// Start nodes
	if utils.MinimalStartUpLogging {
		fmt.Printf("[SERVER]: Starting nodes\n")
	}
	for nodeID := 1; nodeID <= utils.NodeCount; nodeID++ {
		if !utils.MinimalStartUpLogging {
			fmt.Printf("[SERVER]: Creating node %d\n", nodeID)
		}
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
	if utils.MinimalStartUpLogging {
		fmt.Printf("[SERVER]: Sending list of node IP addresses to each node\n")
	}
	for _, nodeAddr := range nodeAddrList {
		if !utils.MinimalStartUpLogging {
			fmt.Printf("[SERVER]: RPC to set neighbors\n")
		}
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

// Find available port on system
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

// Client CLI
func runCLI(client *rpc.Client) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commands (get 'help' to see full options):")

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
				fmt.Printf("Write operation failure: %v\n", err)
			} else {
				fmt.Println("Write operation successful")
			}
		case "forcewrite":
			if argument == "" {
				fmt.Println("Please provide a string to write")
				continue
			}
			var req node.WriteFileRequest
			var res node.WriteFileResponse
			req.Body = argument
			if err := client.Call("Node.ForceWrite", &req, &res); err != nil {
				fmt.Printf("Force write operation failure: %v\n", err)
			} else {
				fmt.Println("Force write operation successful")
			}
		case "read":
			var req node.ReadFileRequest
			var res node.ReadFileResponse
			if err := client.Call("Node.ReadFile", &req, &res); err != nil {
				fmt.Printf("Read operation failure: %v\n", err)
			} else {
				fmt.Println("Data read from file:", res.Data)
			}
		case "info":
			var req node.InfoRequest
			var res node.InfoResponse
			if err := client.Call("Node.Info", &req, &res); err != nil {
				fmt.Printf("Error getting info: %v\n", err)
			} else {
				fmt.Printf("%s\n%s\n", res.AcceptorInfo, res.ProposerInfo)
			}
		case "kill":
			var req node.TerminateRequest
			var res node.TerminateResponse
			if err := client.Call("Node.Terminate", &req, &res); err != nil {
				fmt.Printf("Error calling Terminate RPC method: %v\n", err)
			} else {
				fmt.Println("Termination command sent to all nodes.")
			}
		case "timeout":
			var req node.TimeoutRequest
			var res node.TimeoutResponse
			if err := client.Call("Node.ToggleTimeout", &req, &res); err != nil {
				fmt.Printf("Error toggling timeout: %v\n", err)
			} else {
				fmt.Printf("Timeout ")
				if res.IsTimeout {
					fmt.Printf("on\n")
				} else {
					fmt.Printf("off\n")
				}
			}
		case "stop":
			var req node.StopRequest
			var res node.StopResponse
			if err := client.Call("Node.ToggleStop", &req, &res); err != nil {
				fmt.Printf("Error Stop: %v\n", err)
			} else {
				if res.IsStopped {
					fmt.Printf("Node will stop responding to Paxos requests\n")
				} else {
					fmt.Printf("Node will respond to Paxos requests\n")
				}
			}
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  ping - send ping request to node")
			fmt.Println("  write <string> - write string to file")
			fmt.Println("  read - read string from file")
			fmt.Println("  info - show info about node proposer and acceptor")
			fmt.Println("  help - show this message")
			fmt.Println("  exit - exit program")
		default:
			fmt.Println("Unknown command:", input)
		}
	}
}

// Clear node_data directory
func clearDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
func waitForTermination(client *rpc.Client) {
	var healthCheckReq node.HealthCheckRequest
	var healthCheckRes node.HealthCheckResponse

	for {
		if err := client.Call("Node.HealthCheck", &healthCheckReq, &healthCheckRes); err != nil {
			fmt.Printf("Error checking health status: %v\n", err)
			break
		}

		if healthCheckRes.Status != "OK" {
			fmt.Println("All nodes terminated. Exiting.")
			break
		}

		time.Sleep(1 * time.Second)
	}
}
