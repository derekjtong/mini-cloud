package utils

// Project configs

type NodeConfig struct {
	NodeID    int
	IPAddress string
	Port      int
}

var NodeConfigs = []NodeConfig{
	{NodeID: 1, IPAddress: "127.0.0.1", Port: 0},
	{NodeID: 2, IPAddress: "127.0.0.1", Port: 0},
	{NodeID: 3, IPAddress: "127.0.0.1", Port: 0},
}
