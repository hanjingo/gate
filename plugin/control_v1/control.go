package control_v1

import (
	"time"

	"github.com/hanjingo/container"

	"github.com/hanjingo/gate/com"
	ps "github.com/hanjingo/plugin_system"
	pv4 "github.com/hanjingo/protocol/v4"
)

type agentInfo struct {
	agent     com.AgentI
	id        interface{}    //id
	startTm   time.Time      //开始时间
	subTopics *container.Set //订阅的主题
	apis      *container.Set //可提供的api服务
}
type ControllerV1 struct {
	name                string
	agents              map[interface{}]*agentInfo
	apiAgents           map[uint32]*container.Set //提供api的端点集合 key:api value:端点集合
	subAgents           map[string]*container.Set //订阅了主题的端点集合 key:topic value:端点集合
	codec               *pv4.Codec                //编解码
	newAgentWaitTimeout int                       //新端点等待时长(ms)
	info                *ps.PluginInfo            //插件信息
}

func NewControllerV1() *ControllerV1 {
	back := &ControllerV1{
		name:                PLUGIN_NAME,
		agents:              make(map[interface{}]*agentInfo),
		apiAgents:           make(map[uint32]*container.Set),
		subAgents:           make(map[string]*container.Set),
		codec:               pv4.NewCodec(),
		newAgentWaitTimeout: 3000,
	}
	inf := &ps.PluginInfo{
		Id:          back.name,
		Type:        ps.PLUGIN_TYPE_MEM,
		Version:     "1.0.0",
		Objs:        make(map[string]*ps.Object),
		CallBackMap: make(map[interface{}][]interface{}),
	}
	back.info = inf
	back.reg()
	return back
}

func (c *ControllerV1) Info() *ps.PluginInfo {
	return c.info
}
