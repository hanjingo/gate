package gate

import (
	"github.com/hanjingo/network"
)

type agentV1 struct {
	id   interface{}   //id
	conn network.ConnI //连接
}

func newAgentV1(c network.ConnI) *agentV1 {
	return &agentV1{
		id:   c.GetId(),
		conn: c,
	}
}

func (a *agentV1) GetId() interface{} {
	return a.id
}

func (a *agentV1) Write(args ...[]byte) error {
	return nil
}

func (a *agentV1) Read() ([]byte, error) {
	return nil, nil
}

func (a *agentV1) Close() {}
