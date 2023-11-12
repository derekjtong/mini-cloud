package main

import (
	"log"
	"net/rpc"
)

// rpc client

type Args struct{}

func main() {

	hostname := "localhost"
	port := ":1122"

	var reply string

	args := Args{}

	client, err := rpc.DialHTTP("tcp", hostname+port)
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	// Call normally takes service name.function name, args and
	// the address of the variable that hold the reply. Here we
	// have no args in the demo therefore we can pass the empty
	// args struct.
	err = client.Call("BannerMessageServer.GetBannerMessage", args, &reply)
	if err != nil {
		log.Fatal("error", err)
	}

	// log the result
	log.Printf("%s\n", reply)
}
