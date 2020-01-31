package stream_v1

import "github.com/hanjingo/gate/com"

//操作码
const (
	OP_LIMIT_STREAM           = com.MASK_CONTROL | 0x9  //控流
	OP_SET_LIMIT_SPEED_CONF   = com.MASK_CONTROL | 0x10 //设置流速
	OP_SET_LIMIT_SCANDUR_CONF = com.MASK_CONTROL | 0x11 //设置扫描间隔
)
