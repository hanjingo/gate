package control_v1

import (
	"github.com/hanjingo/container"

	ps "github.com/hanjingo/plugin_system"
	pv4 "github.com/hanjingo/protocol/v4"
)

const NAME = "ControllerV1" //插件名字
const VERSION = "1.0.0"     //插件版本

type ControllerV1 struct {
	name                string
	apiAgents           map[uint32]*container.Set //提供api的端点集合 key:api value:端点集合
	subAgents           map[string]*container.Set //订阅了主题的端点集合 key:topic value:端点集合
	codec               *pv4.Codec                //编解码
	newAgentWaitTimeout int                       //新端点等待时长(ms)
	info                *ps.PluginInfo            //插件信息
}

func New() ps.PluginI {
	back := &ControllerV1{
		name:                NAME,
		apiAgents:           make(map[uint32]*container.Set),
		subAgents:           make(map[string]*container.Set),
		codec:               pv4.NewCodec(),
		newAgentWaitTimeout: 3000,
	}
	inf := &ps.PluginInfo{
		Id:          back.name,
		Type:        ps.PLUGIN_TYPE_MEM,
		Version:     VERSION,
		Objs:        make(map[string]*ps.Object),
		CallBackMap: make(map[interface{}]interface{}),
	}
	back.info = inf
	back.reg()
	return back
}

func (c *ControllerV1) Info() *ps.PluginInfo {
	return c.info
}

func (c *ControllerV1) agents() map[interface{}]*agentInfo {
	return agents
}
