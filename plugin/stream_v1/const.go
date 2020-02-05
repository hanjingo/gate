package stream_v1

import (
	"github.com/hanjingo/logger"
)

var log = logger.GetDefaultLogger()

//权限掩码
const (
	MASK_CONTROL uint32 = 0x80000000 //控制网关权限
)

//操作码
const (
	OP_LIMIT_STREAM           uint32 = MASK_CONTROL | 0x9  //控流
	OP_SET_LIMIT_SPEED_CONF   uint32 = MASK_CONTROL | 0x10 //设置流速
	OP_SET_LIMIT_SCANDUR_CONF uint32 = MASK_CONTROL | 0x11 //设置扫描间隔
)
