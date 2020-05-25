package filt_v1

import (
	ps "github.com/hanjingo/golib/plugin_system"
)

const NAME = "FiltV1"   //插件名字
const VERSION = "1.0.0" //插件版本

type FiltV1 struct {
	name string
	info *ps.PluginInfo
}

func New() ps.PluginI {
	back := &FiltV1{
		name: NAME,
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

func (f *FiltV1) Info() *ps.PluginInfo {
	return f.info
}

func (f *FiltV1) agents() map[interface{}]*agentInfo {
	return agents
}
