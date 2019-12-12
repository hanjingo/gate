package com

import (
	"time"

	"github.com/hanjingo/util"
)

type Msg struct {
	OpCode    uint32
	Receivers []uint64
	Sender    uint64
	TimeStamp string
	Content   interface{}
}

func NewMsg() *Msg {
	return &Msg{
		TimeStamp: util.TimeToTimeStamp(time.Now()),
	}
}
