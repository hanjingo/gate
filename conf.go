package gate

import (
	"github.com/hanjingo/network"
)

type GateConfig struct {
	Id           uint8
	ConfServAddr string
	UserName     string
	PassWord     string
	Plugins      map[string]*PluginConfig
	Servers      map[string]*network.ServerConfig
}

type PluginConfig struct {
	Name string
}

func NewGateConfig() *GateConfig {
	back := &GateConfig{
		Plugins: make(map[string]*PluginConfig),
		Servers: make(map[string]*network.ServerConfig),
	}
	return back
}
