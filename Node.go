package main
import (  
	"fmt"  
	"log"  
	"net"  
)
type RPCRequest struct {  
	NodeID   int  
	Operation string  
	Data     []byte  
}
type RPCResponse struct {  
	Status   bool  
	Message  string  
	Data     []byte  
}
func main() {  
	// Establish connections to all nodes  
	connections := make([]*net.TCPConn, 3)  
	for i := 0; i < 3; i++ {  
		conn, err := net.Dial("tcp", "127.0.0.1:5555" + fmt.Sprintf("%d", i))  
		if err != nil {  
			log.Fatal(err)  
		}  
		connections[i] = conn  
	}
	// Simulate file creation  
	createFile(connections)  
