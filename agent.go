package gate

import (
	net "github.com/hanjingo/golib/network"
	util "github.com/hanjingo/golib/util"
)

var gen = util.GetUuidGenerator()

type agentV1 struct {
	id      interface{}  //id
	conn    net.SessionI //连接
	isValid bool         //是否可用
}

func newAgentV1(c net.SessionI) *agentV1 {
	id, err := gen.GenerateUuidUint64()
	if err != nil {
		return nil
	}
	back := &agentV1{
		id:      id,
		conn:    c,
		isValid: false, //经过验证才让他可用
	}
	return back
}

func (a *agentV1) GetId() interface{} {
	return a.id
}

func (a *agentV1) Write(args ...[]byte) error {
	return a.conn.WriteMsg(args...)
}

func (a *agentV1) Read() ([]byte, error) {
	return a.conn.ReadMsg()
}

func (a *agentV1) Close() {
	a.conn.Close()
}

func (a *agentV1) IsValid() bool {
	return a.isValid
}
