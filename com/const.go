package com

import (
	"github.com/hanjingo/util"
)

const START_TIME_STR string = "2019-01-01 00:00:00"

var START_TIME = util.TimeStampToTime(START_TIME_STR)

//未知编码id
const UNKNOWN_CODECID uint8 = 0

//一些系统消息
const (
	AGENT_CONNECT uint32 = 0xa0000000 | 0x1 //建立端点
	AGENT_CLOSE   uint32 = 0xa0000000 | 0x2 //关闭端点
)

