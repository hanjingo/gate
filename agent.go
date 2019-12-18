package gate

import (
	"github.com/hanjingo/network"
)

type agentV1 struct {
	id      interface{}   //id
	conn    network.ConnI //连接
	isValid bool          //是否可用
}

func newAgentV1(c network.ConnI) *agentV1 {
	back := &agentV1{
		id:      c.GetId(),
		conn:    c,
		isValid: false, //经过验证才让他可用
	}
	back.conn.Run()
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
