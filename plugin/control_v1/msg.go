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
