package gate

import (
	"net/http"
	"strings"
	"sync"

	com "github.com/hanjingo/gate/com"
	gplugin "github.com/hanjingo/gate/plugin"
	logger "github.com/hanjingo/golib/logger"
	network "github.com/hanjingo/golib/network"
	plugin "github.com/hanjingo/golib/plugin"
	pv4 "github.com/hanjingo/golib/protocol/v4"
)

var log = logger.GetDefaultLogger()

type Gate struct {
	wg      *sync.WaitGroup
	conf    *GateConfig                      //配置
	agents  map[interface{}]com.AgentI       //端点集合
	conns   map[network.SessionI]interface{} //连接和端点映射
	servers map[interface{}]network.ServerI  //tcp/ws/http服务器集合
	hubs    *plugin.Hubs                     //插件系统
	codec   *pv4.Codec                       //编解码器
}

func NewGate(conf *GateConfig) *Gate {
	back := &Gate{
		wg:      new(sync.WaitGroup),
		conf:    conf,
		agents:  make(map[interface{}]com.AgentI),
		conns:   make(map[network.SessionI]interface{}),
		servers: make(map[interface{}]network.ServerI),
		hubs:    plugin.NewHubs(),
		codec:   pv4.NewCodec(),
	}
	back.reg()
	//加载插件
	for _, e := range conf.Plugins {
		if f, ok := gplugin.GetPlugins()[e.Name]; ok {
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
	gate.conns[c] = agent.GetId()
	if gate.hubs != nil {
		f := func(id interface{}, hub plugin.HubI) {
			hub.Call(com.AGENT_CONNECT, agent)
		}
		gate.hubs.RangeHub(f)
	}
}

//处理连接断开
func (gate *Gate) onConnClose(c network.SessionI) {
	agent, ok := gate.agents[c]
	if !ok {
		return
	}
	delete(gate.agents, agent.GetId())
	delete(gate.conns, c)
	if gate.hubs != nil {
		f := func(id interface{}, hub plugin.HubI) {
			hub.Call(com.AGENT_CLOSE, agent)
		}
		gate.hubs.RangeHub(f)
	}
}

//处理收到的消息
func (gate *Gate) handleMsg(c network.SessionI, data []byte) {
	id, ok := gate.conns[c]
	if !ok {
		return
	}
	agent, ok := gate.agents[id]
	if !ok {
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
