package gate

import (
	"github.com/hanjingo/network"
)

type PluginConfig struct {
	Name string
}

type GateConfig struct {
	Id      uint8
	Servers map[string]*network.ServerConfig
	Plugins map[string]*PluginConfig //key:plugin名字 value:plugin设置
}

func NewGateConfig() *GateConfig {
	back := &GateConfig{
		Servers: make(map[string]*network.ServerConfig),
		Plugins: make(map[string]*PluginConfig),
	}
	return back
}
