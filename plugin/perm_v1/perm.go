package perm_v1

import (
	ps "github.com/hanjingo/plugin_system"
	pv4 "github.com/hanjingo/protocol/v4"
)

const NAME = "PermV1"   //插件名字
const VERSION = "1.0.0" //插件版本

type PermV1 struct {
	name  string
	codec *pv4.Codec
	info  *ps.PluginInfo
}

func New() ps.PluginI {
	back := &PermV1{
		name:  NAME,
		codec: pv4.NewCodec(),
	}
	back.info = &ps.PluginInfo{
		Id:          back.name,
		Type:        ps.PLUGIN_TYPE_MEM,
		Version:     VERSION,
		Objs:        make(map[string]*ps.Object),
		CallBackMap: make(map[interface{}]interface{}),
	}
	back.reg()
	return back
}

func (p *PermV1) Info() *ps.PluginInfo {
	return p.info
}

func (p *PermV1) agents() map[interface{}]*agentInfo {
	return agents
}
