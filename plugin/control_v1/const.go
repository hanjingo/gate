package control_v1

import "github.com/hanjingo/gate/com"

//操作码
const (
	OP_PROTECT = 0xff //保护掩码

	//基础功能
	OP_ROUTE         = com.MASK_NORM | 0x1 //路由
	OP_PING          = com.MASK_NORM | 0x2 //ping
	OP_PONG          = com.MASK_NORM | 0x3 //pong
	OP_NEW_AGENT_REQ = com.MASK_NORM | 0x4 //建立端点
	OP_NEW_AGENT_RSP = com.MASK_NORM | 0x5 //建立端点返回

	//客户端
	OP_SUB   = com.MASK_CLI | 0x1 //订阅
	OP_UNSUB = com.MASK_CLI | 0x2 //取消订阅

	//服务端
	OP_REG_SERVER   = com.MASK_SERVER | 0x1 //注册服务
	OP_UNREG_SERVER = com.MASK_SERVER | 0x2 //取消注册
	OP_MULTCAST     = com.MASK_SERVER | 0x3 //组播
	OP_BROADCAST    = com.MASK_SERVER | 0x4 //广播
	OP_PUB          = com.MASK_SERVER | 0x5 //发布

	//控制码 (1 ~ 63)
	OP_CLOSE_AGENT = com.MASK_CONTROL | 0x1 //关闭节点

	OP_SET_FILT_AGENT   = com.MASK_CONTROL | 0x2 //设置过滤
	OP_SET_UNFILT_AGENT = com.MASK_CONTROL | 0x3 //解除过滤
	OP_ON_FILT          = com.MASK_CONTROL | 0x4 //过滤消息

	OP_SET_AGENT_LIMIT_STREAM   = com.MASK_CONTROL | 0x5 //设置节点控流限制
	OP_SET_AGENT_UNLIMIT_STREAM = com.MASK_CONTROL | 0x6 //设置节点不控流
	OP_ON_STREAM                = com.MASK_CONTROL | 0x7 //过滤消息

	OP_CHANGE_PERM = com.MASK_CONTROL | 0x8 //变更权限

	OP_LIMIT_STREAM = com.MASK_CONTROL | 0x9 //限制流量
)
