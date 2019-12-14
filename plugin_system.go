package gate

import (
	"github.com/hanjingo/gate/com"
)

type PluginSystem struct {
	plugins map[string]com.PluginI
	//taskTree 任务树
}

func newPluginSystem() *PluginSystem {
	return &PluginSystem{
		plugins: make(map[string]com.PluginI),
	}
}

func (ps *PluginSystem) onMsg(agent com.AgentI, data []byte) error {
	return nil
}

func (ps *PluginSystem) onNewAgent(agent com.AgentI) error {
	//todo
	return nil
}

func (ps *PluginSystem) onAgentClose(agent com.AgentI) error {
	//todo
	return nil
}

func (ps *PluginSystem) addPlugin(plugin com.PluginI) {
	//todo
}

func (ps *PluginSystem) delPlugin(name string) {
	//todo
}
