package gate

import (
	"net/http"
	"strings"
	"sync"

	"github.com/hanjingo/gate/plugin"
	"github.com/hanjingo/logger"

	"github.com/hanjingo/gate/com"
	"github.com/hanjingo/network"
	ps "github.com/hanjingo/plugin_system"
	pv4 "github.com/hanjingo/protocol/v4"
)

var log = logger.GetDefaultLogger()

type Gate struct {
	wg      *sync.WaitGroup
	conf    *GateConfig                     //配置
	agents  map[interface{}]com.AgentI      //端点集合
	servers map[interface{}]network.ServerI //tcp/ws/http服务器集合
	hubs    *ps.Hubs                        //插件系统
	codec   *pv4.Codec                      //编解码器
}

func NewGate(conf *GateConfig) *Gate {
	back := &Gate{
		wg:      new(sync.WaitGroup),
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
			p := f()
			if err := back.hubs.LoadPlugin(p); err != nil {
				panic(err)
			}
			log.Info("插件:%v 被加载", p.Info().Id)
		}
	}
	//加载服务器
	for _, sconf := range conf.Servers {
		var s network.ServerI
		var err error
		switch strings.ToUpper(sconf.Type) {
		case "TCP":
			s, err = network.NewTcpServer(sconf, back.onConnClose,
				back.onNewConn, back.handleMsg)
		case "WS":
			s, err = network.NewWsServer(sconf, back.onConnClose,
				back.onNewConn, back.handleMsg)
		}
		if err != nil {
			log.Error("服务器启动失败,错误:%v", err)
			continue
		}
		back.servers[sconf.Name] = s
	}
	return back
}

//跑起来
func (gate *Gate) Run() {
	//监听cmd服务器
	go func() {
		if err := http.ListenAndServe(gate.conf.ConfServAddr, nil); err != nil {
			panic(err)
		}
		log.Info("cmd服务器已启动...")
	}()

	//监听tcp/ws
	for name, s := range gate.servers {
		s.Run(gate.wg)
		log.Info("服务器:%v 已启动", name)
	}
}

//等待
func (gate *Gate) Wait() {
	gate.wg.Wait()
}

//处理新建连接
func (gate *Gate) onNewConn(c network.SessionI) {
	agent := newAgentV1(c)
	gate.agents[agent.GetId()] = agent
	if gate.hubs != nil {
		f := func(id interface{}, hub ps.HubI) {
			hub.Call(com.AGENT_CONNECT, agent)
		}
		gate.hubs.RangeHub(f)
	}
	log.Info("连接:%v被建立", agent.GetId())
}

//处理连接断开
func (gate *Gate) onConnClose(c network.SessionI) {
	log.Debug("收到连接:%v被摧毁通知", c.GetId())
	agent, ok := gate.agents[c.GetId()]
	if !ok {
		return
	}
	delete(gate.agents, c.GetId())
	if gate.hubs != nil {
		f := func(id interface{}, hub ps.HubI) {
			hub.Call(com.AGENT_CLOSE, agent)
		}
		gate.hubs.RangeHub(f)
	}
	log.Info("连接:%v被摧毁", c.GetId())
}

//处理收到的消息
func (gate *Gate) handleMsg(agentId uint64, data []byte) {
	agent, ok := gate.agents[agentId]
	if !ok {
		log.Error("收到未知端点:%v的消息", agentId)
		return
	}
	//解析操作码
	opcode, err := gate.codec.ParseOpCode(data)
	if err != nil {
		return
	}
	//捕获系统消息
	if opcode&com.MASK_SYS == com.MASK_SYS {
		//todo
		switch opcode {
		case com.AGENT_CONNECT:
			log.Info("捕获系统消息: AGENT_CONNECT")
		case com.AGENT_CLOSE:
			log.Info("捕获系统消息: AGENT_CLOSE")
		default:
			log.Error("捕获未知系统消息:%b", opcode)
		}
		return
	}
	//发给插件
	if gate.hubs != nil {
		gate.hubs.Call(opcode, agent, data)
	}
}

//返回网关信息
func (gate *Gate) info() *GateInfo {
	//todo
	return nil
}
