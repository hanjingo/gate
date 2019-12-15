package gate

import (
	"github.com/hanjingo/gate/com"
)

//插件系统
type PluginSystem struct {
	plugins map[string]com.PluginI
	tasks   []string
}

func newPluginSystem() *PluginSystem {
	back := &PluginSystem{
		plugins: make(map[string]com.PluginI),
		tasks:   []string{},
	}
	return back
}

//处理消息
func (ps *PluginSystem) onMsg(agent com.AgentI, data []byte) error {
	go func() {
		var err error
		for _, name := range ps.tasks {
			if p, ok := ps.plugins[name]; ok {
				data, err = p.OnMsg(agent, data)
				if err != nil {
					return
				}
			}
		}
	}()
	return nil
}

//建立端点
func (ps *PluginSystem) onNewAgent(agent com.AgentI) error {
	for _, p := range ps.plugins {
		if err := p.OnNewAgent(agent); err != nil {
			return err
		}
	}
	return nil
}

//关闭端点
func (ps *PluginSystem) onAgentClose(agent com.AgentI) error {
	for _, p := range ps.plugins {
		p.OnAgentClose(agent)
	}
	return nil
}

//添加插件
func (ps *PluginSystem) addPlugin(plugin com.PluginI) {
	ps.plugins[plugin.Name()] = plugin
	ps.tasks = append(ps.tasks, plugin.Name())
}

//删除插件
func (ps *PluginSystem) delPlugin(name string) {
	for i, pname := range ps.tasks {
		if pname == name {
			ps.tasks = append(ps.tasks[:i], ps.tasks[i+1:]...)
		}
	}
	delete(ps.plugins, name)
}
