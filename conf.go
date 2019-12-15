package gate

import (
	"github.com/hanjingo/network"
)

type GateConfig struct {
	Id      uint8
	Servers map[string]*network.ServerConfig
}

func NewGateConfig() *GateConfig {
	back := &GateConfig{
		Servers: make(map[string]*network.ServerConfig),
	}
	return back
}
