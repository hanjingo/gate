package perm_v1

import (
	"github.com/hanjingo/logger"
)

var log = logger.GetDefaultLogger()

//权限掩码
const (
	MASK_NORM    uint32 = 0x10000000 //普通权限
	MASK_CONTROL uint32 = 0x80000000 //控制网关权限
)

//操作码
const (
	OP_CHANGE_PERM uint32 = MASK_CONTROL | 0x8 //变更权限
)
