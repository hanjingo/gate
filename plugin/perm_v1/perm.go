package perm_v1

import (
	plugin "github.com/hanjingo/golib/plugin"
	pv4 "github.com/hanjingo/golib/protocol/v4"
	types "github.com/hanjingo/golib/types"
)

const NAME = "PermV1"   //插件名字
const VERSION = "1.0.0" //插件版本

type PermV1 struct {
	name  string
	codec *pv4.Codec
	info  *plugin.PluginInfo
}

func New() plugin.PluginI {
	back := &PermV1{
		name:  NAME,
		codec: pv4.NewCodec(),
	}
	back.info = &plugin.PluginInfo{
		Id:          back.name,
		Type:        plugin.PTYPE_MEM,
		Version:     VERSION,
		Objs:        make(map[string]*types.Object),
		CallBackMap: make(map[interface{}]interface{}),
	}
	back.reg()
	return back
}

func (p *PermV1) Info() *plugin.PluginInfo {
	return p.info
}

func (p *PermV1) agents() map[interface{}]*agentInfo {
	return agents
}
