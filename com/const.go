package com

import (
	"github.com/hanjingo/util"
)

const START_TIME_STR string = "2019-01-01 00:00:00"

var START_TIME = util.TimeStampToTime(START_TIME_STR)

//未知编码id
const UNKNOWN_CODECID uint8 = 0

//权限掩码
const (
	MASK         = 0x1fffffff //掩码
	MASK_ZERO    = 0x0        //0权限
	MASK_NORM    = 0x20000000 //普通权限
	MASK_CLI     = 0x40000000 //客户端权限
	MASK_SERVER  = 0x60000000 //服务端权限
	MASK_CONTROL = 0x80000000 //控制网关权限
	MASK_SYS     = 0xa0000000 //系统权限
)

//一些系统消息
const (
	AGENT_CONNECT uint32 = MASK_SYS | 0x1 //建立端点
	AGENT_CLOSE   uint32 = MASK_SYS | 0x2 //关闭端点
)
