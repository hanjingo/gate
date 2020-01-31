package filt_v1

import "github.com/hanjingo/gate/com"

//消息id (64~128)
const (
	OP_SET_FILT_AGENT   = com.MASK_CONTROL | 0x2 //设置过滤
	OP_SET_UNFILT_AGENT = com.MASK_CONTROL | 0x3 //解除过滤
	OP_ON_FILT          = com.MASK_CONTROL | 0x4 //过滤消息
)
