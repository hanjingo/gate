package filt_v1

import (
	plugin "github.com/hanjingo/golib/plugin"
	types "github.com/hanjingo/golib/types"
)

const NAME = "FiltV1"   //插件名字
const VERSION = "1.0.0" //插件版本

type FiltV1 struct {
	name string
	info *plugin.PluginInfo
}

func New() plugin.PluginI {
	back := &FiltV1{
		name: NAME,
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

func (f *FiltV1) Info() *plugin.PluginInfo {
	return f.info
}

func (f *FiltV1) agents() map[interface{}]*agentInfo {
	return agents
}
