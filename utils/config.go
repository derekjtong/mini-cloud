package utils

// Project configs

type NodeConfig struct {
	NodeID    int
	IPAddress string
	Port      int
}

var NodeConfigs = []NodeConfig{
	{NodeID: 1, IPAddress: "127.0.0.1", Port: 65432},
	{NodeID: 2, IPAddress: "127.0.0.1", Port: 65431},
	{NodeID: 3, IPAddress: "127.0.0.1", Port: 65430},
}
