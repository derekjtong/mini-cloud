package main

import (
    "log"
    "net"
    "net/http"
    "net/rpc"
)

// an RPC server in Go

type Args struct{}

type BannerMessageServer string

func (t *BannerMessageServer) GetBannerMessage(args *Args, reply *string) error {
    *reply = "This is a message from the RPC server"
    return nil
}

func main() {

    // create and register the rpc
    banner := new(BannerMessageServer)
    rpc.Register(banner)
    rpc.HandleHTTP()

    // set a port for the server
    port := ":1122"

    // listen for requests on 1122
    listener, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatal("listen error: ", err)
    }

    http.Serve(listener, nil)
}