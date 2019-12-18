package com

type PluginI interface {
	Name() string
	Set(string, interface{}) //设置一些参数和函数
	OnNewAgent(AgentI) error
	OnAgentClose(AgentI) error
	OnMsg(AgentI, []byte) ([]byte, error)
}

type AgentI interface {
	IsValid() bool
	GetId() interface{}
	Write(args ...[]byte) error
	Read() ([]byte, error)
	Close()
}
