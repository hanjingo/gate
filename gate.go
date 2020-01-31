package gate

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/hanjingo/gate/plugin"

	"github.com/hanjingo/gate/com"
	"github.com/hanjingo/network"
	ps "github.com/hanjingo/plugin_system"
	pv4 "github.com/hanjingo/protocol/v4"
)

type Gate struct {
	conf    *GateConfig                     //配置
	agents  map[interface{}]com.AgentI      //端点集合
	servers map[interface{}]network.ServerI //tcp/ws/http服务器集合
	hubs    *ps.Hubs                        //插件系统
	codec   *pv4.Codec                      //编解码器
}

func NewGate(conf *GateConfig) *Gate {
	back := &Gate{
		conf:    conf,
		agents:  make(map[interface{}]com.AgentI),
		servers: make(map[interface{}]network.ServerI),
		hubs:    ps.NewHubs(),
		codec:   pv4.NewCodec(),
	}
	back.reg()
	//加载插件
	for _, e := range conf.Plugins {
		if f, ok := plugin.GetPlugins()[e.Name]; ok {
			fmt.Println("插件:", e.Name, "被加载")
			if err := back.hubs.LoadPlugin(f()); err != nil {
				panic(err)
			}
		}
	}
	return back
}

//跑起来
func (gate *Gate) Run(wg *sync.WaitGroup) {
	//监听http
	http.ListenAndServe(gate.conf.ConfServAddr, nil)
}

//处理新建连接
func (gate *Gate) onNewConn(c network.ConnI) {
	agent := newAgentV1(c)
	if gate.hubs != nil {
		for _, hub := range gate.hubs.GetHubs() {
			hub.Call(com.AGENT_CONNECT, agent)
		}
	}
	gate.agents[agent.GetId()] = agent
}

//处理连接断开
func (gate *Gate) onConnClose(c network.ConnI) {
	agent, ok := gate.agents[c.GetId()]
	if !ok {
		return
	}
	if gate.hubs != nil {
		for _, hub := range gate.hubs.GetHubs() {
			hub.Call(com.AGENT_CLOSE, agent)
		}
	}
	delete(gate.agents, c.GetId())
}

//处理收到的消息
func (gate *Gate) handleMsg(agentId uint64, data []byte) {
	agent, ok := gate.agents[agentId]
	if !ok {
		return
	}
	//解析操作码
	opcode, err := gate.codec.ParseOpCode(data)
	if err != nil {
		return
	}
	//屏蔽系统消息
	if opcode&com.MASK_SYS != 0x0 {
		return
	}
	//发给插件
	if gate.hubs != nil {
		gate.hubs.Call(opcode, agent)
	}
}

//返回网关信息
func (gate *Gate) info() *GateInfo {
	//todo
	return nil
}
