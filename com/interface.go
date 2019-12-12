package com

type PluginI interface {
	Name() string
	OnNewAgent(AgentI) error
	OnAgentClose(AgentI) error
	OnMsg(*Msg) (*Msg, error)
}

type AgentI interface {
	GetId() interface{}
	Write(args ...[]byte) error
	Read() ([]byte, error)
	Close()
}
