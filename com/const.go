package com

import (
	"github.com/hanjingo/util"
)

const START_TIME_STR string = "2019-01-01 00:00:00"

var START_TIME = util.TimeStampToTime(START_TIME_STR)

//未知编码id
const UNKNOWN_CODECID uint8 = 0

//权限掩码
const MASK_SYS = 0xa0000000 //系统权限

//一些系统消息
const (
	AGENT_CONNECT uint32 = MASK_SYS | 0x1 //建立端点
	AGENT_CLOSE   uint32 = MASK_SYS | 0x2 //关闭端点
)
