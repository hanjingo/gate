package golang

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

//操作码
const (
	OP_PROTECT = 0xff //保护掩码

	//基础功能
	OP_ROUTE         = MASK_NORM | 0x1 //路由
	OP_PING          = MASK_NORM | 0x2 //ping
	OP_PONG          = MASK_NORM | 0x3 //pong
	OP_NEW_AGENT     = MASK_NORM | 0x4 //建立端点
	OP_NEW_AGENT_RSP = MASK_NORM | 0x5 //建立端点返回

	//客户端
	OP_SUB   = MASK_CLI | 0x1 //订阅
	OP_UNSUB = MASK_CLI | 0x2 //取消订阅

	//服务端
	OP_REG_SERVER   = MASK_SERVER | 0x1 //注册服务
	OP_UNREG_SERVER = MASK_SERVER | 0x2 //取消注册
	OP_MULTCAST     = MASK_SERVER | 0x3 //组播
	OP_BROADCAST    = MASK_SERVER | 0x4 //广播
	OP_PUB          = MASK_SERVER | 0x5 //发布

	//控制码
	OP_CLOSE_AGENT = MASK_CONTROL | 0x1 //关闭节点

	OP_SET_FILT_AGENT   = MASK_CONTROL | 0x2 //设置过滤
	OP_SET_UNFILT_AGENT = MASK_CONTROL | 0x3 //解除过滤
	OP_ON_FILT          = MASK_CONTROL | 0x4 //过滤消息

	OP_SET_AGENT_LIMIT_STREAM   = MASK_CONTROL | 0x5 //设置节点控流限制
	OP_SET_AGENT_UNLIMIT_STREAM = MASK_CONTROL | 0x6 //设置节点不控流
	OP_ON_STREAM                = MASK_CONTROL | 0x7 //控流消息

	OP_CHANGE_PERM = MASK_CONTROL | 0x8 //变更权限
)
