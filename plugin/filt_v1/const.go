package filt_v1

import (
	"github.com/hanjingo/logger"
)

var log = logger.GetDefaultLogger()

//掩码
const MASK_CONTROL uint32 = 0x80000000 //控制网关权限

//消息id (64~128)
const (
	OP_SET_FILT_AGENT   uint32 = MASK_CONTROL | 0x2 //设置过滤
	OP_SET_UNFILT_AGENT uint32 = MASK_CONTROL | 0x3 //解除过滤
	OP_ON_FILT          uint32 = MASK_CONTROL | 0x4 //过滤消息
)
