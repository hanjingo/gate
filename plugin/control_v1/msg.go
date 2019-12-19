package control_v1

type MsgNewAgentReq struct {
	Token string
}
type MsgNewAgentRsp struct {
	Result bool
	Id     uint64
}

type MsgPing struct {
}

type MsgPong struct {
}

type MsgSub struct {
	Topics []string
}

type MsgUnSub struct {
	Topics []string
}

type MsgPub struct {
	Topic string
	Info  interface{}
}

type MsgRegApi struct {
	Apis []uint32
	Id   uint64
}

type MsgUnRegApi struct {
	Apis []uint32
	Id   uint64
}

type MsgControl struct {
	Cmd uint32
	Args []interface{}
}
