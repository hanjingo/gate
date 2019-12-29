package com

type AgentI interface {
	IsValid() bool
	GetId() interface{}
	Write(args ...[]byte) error
	Read() ([]byte, error)
	Close()
}
