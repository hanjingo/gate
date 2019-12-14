package com

type PluginI interface {
	Name() string
	OnNewAgent(AgentI) error
	OnAgentClose(AgentI) error
	OnMsg(AgentI, []byte) ([]byte, error)
}

type AgentI interface {
	GetId() interface{}
	Write(args ...[]byte) error
	Read() ([]byte, error)
	Close()
}
