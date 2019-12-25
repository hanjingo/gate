package plugin

import (
	"github.com/hanjingo/gate/com"
	ps "github.com/hanjingo/plugin_system"
	pv4 "github.com/hanjingo/protocol/v4"
)

var Ps *PluginSystem

func GetPluginSystem() *PluginSystem {
	if Ps == nil {
		Ps = &PluginSystem{
			hubs:  ps.NewHubs(),
			codec: pv4.NewCodec(),
		}
	}
	return Ps
}

//插件系统
type PluginSystem struct {
	hubs  *ps.Hubs
	codec *pv4.Codec //编解码器
}

//加载插件
func (ps *PluginSystem) LoadPlugin(addr interface{}) error {
	return nil
}

//卸载插件
func (ps *PluginSystem) UnLoadPlugin(id interface{}) error {
	return ps.hubs.UnLoadPlugin(id)
}

//新建端点
func (ps *PluginSystem) OnNewAgent(agent com.AgentI) error {
	return nil
}

//关闭端点
func (ps *PluginSystem) OnAgentClose(agent com.AgentI) error {
	return nil
}

//处理消息
func (ps *PluginSystem) OnMsg(agent com.AgentI, data []byte) error {
	opcode, err := ps.codec.ParseOpCode(data)
	if err != nil {
		return err
	}
	ps.hubs.Call(opcode, data)
	return nil
}
