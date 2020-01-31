package perm_v1

import "github.com/hanjingo/gate/com"

//操作码
const (
	OP_CHANGE_PERM = com.MASK_CONTROL | 0x8 //变更权限
)
