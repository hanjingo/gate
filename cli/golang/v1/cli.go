package v1

import (
	"errors"
	"fmt"
	"reflect"

	ws "github.com/gorilla/websocket"
	"github.com/hanjingo/network"
	"github.com/hanjingo/protocol"
	pv4 "github.com/hanjingo/protocol/v4"
)

//网关客户端
type GateCli struct {
	id      uint64
	conn    network.SessionI
	codec   protocol.CodecI
	callMap map[interface{}]*Handler
	onClose func(*GateCli)
}

func NewGateCli() *GateCli {
	back := &GateCli{
		codec:   pv4.NewCodec(),
		callMap: make(map[interface{}]*Handler),
	}
	return back
}

//拨号
func (cli *GateCli) Dial(dialType string, addr string, token string, conf *network.SessionConfig) error {
	var conn network.SessionI
	var err error
	c, _, err := ws.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	conn, err = network.NewWsConn(conf, c, cli.onConnClose, nil)
	if err != nil {
		return err
	}

	if err := readMsg(conn, cli.codec, &MsgAgentConnSucc{}); err != nil {
		return err
	}

	//发请求
	if err := writeMsg(conn, cli.codec, OP_NEW_AGENT_REQ, 0, &MsgNewAgentReq{Token: token}, nil); err != nil {
		return err
	}
	//收回复
	rsp := &MsgNewAgentRsp{}
	if err := readMsg(conn, cli.codec, rsp); err != nil {
		return err
	}
	if !rsp.Result {
		return errors.New("网关拒绝建立连接")
	}
	cli.conn = conn
	cli.id = rsp.Id
	fmt.Println("与网关:", addr, "连接建立")
	return nil
}

//注册回调
func (cli *GateCli) RegHandler(id interface{}, f interface{}, arg ...interface{}) error {
	if _, ok := cli.callMap[id]; ok {
		return errors.New("已经注册，无需再次注册")
	}
	inv := []reflect.Value{}
	fun := func(in ...interface{}) {
		for _, e := range in {
			inv = append(inv, reflect.ValueOf(e))
		}
		fv := reflect.ValueOf(f)
		fv.Call(inv)
	}
	h := &Handler{Id: id, Call: fun, Arg: arg}
	cli.callMap[id] = h
	return nil
}

//处理关闭
func (cli *GateCli) onConnClose(s network.SessionI) {
	if cli.onClose != nil {
		cli.onClose(cli)
	}
}

//回调
func (cli *GateCli) OnMsg(session *network.SessionI, data []byte) {
	codec := cli.codec.(*pv4.Codec)
	id, err := codec.ParseOpCode(data)
	if err != nil {
		return
	}
	sender, err := codec.ParseSender(data)
	if err != nil {
		return
	}
	h, ok := cli.callMap[id]
	if !ok {
		return
	}
	var args []interface{}
	args = append(args, cli)//cli
	args = append(args, sender)//agentid
	t := reflect.ValueOf(h.Arg).Type()
	v := reflect.New(t)
	if v.CanInterface() {
		temp := v.Interface()
		if err := codec.UnFormat(data, temp); err != nil {
			return
		}
		args = append(args, temp)
	} else {
		args = append(args, v)
	}
	h.Call(args...) //call(cli, agentid, req)
}

//路由
func (cli *GateCli) Route(msg interface{}, recvs ...uint64) error {
	return writeMsg(cli.conn, cli.codec, OP_ROUTE, cli.id, msg, recvs)
}

//ping
func (cli *GateCli) Ping() error {
	return writeMsg(cli.conn, cli.codec, OP_PING, cli.id,
		&MsgPing{}, nil)
}

//请求
func (cli *GateCli) Req(opCode uint32, msg interface{}, recv ...uint64) error {
	return writeMsg(cli.conn, cli.codec, MASK_CLI&opCode, cli.id, msg, recv)
}

//返回
func (cli *GateCli) Rsp(opCode uint32, msg interface{}, recv ...uint64) error {
	return writeMsg(cli.conn, cli.codec, MASK_CLI&opCode, cli.id, msg, recv)
}

//订阅
func (cli *GateCli) Sub(topics ...string) error {
	return writeMsg(cli.conn, cli.codec, OP_SUB, cli.id,
		&MsgSub{Topics: topics}, nil)
}

//发布
func (cli *GateCli) Pub() {

}

//取消订阅
func (cli *GateCli) UnSub(topics ...string) error {
	return writeMsg(cli.conn, cli.codec, OP_UNSUB, cli.id,
		&MsgUnSub{Topics: topics}, nil)
}

//注册服务
func (cli *GateCli) RegApi(apis ...uint32) error {
	return writeMsg(cli.conn, cli.codec, OP_REG_SERVER, cli.id,
		&MsgRegApi{Apis: apis, Id: cli.id}, nil)
}

//取消注册
func (cli *GateCli) UnRegApi(apis ...uint32) error {
	return writeMsg(cli.conn, cli.codec, OP_UNREG_SERVER, cli.id,
		&MsgUnRegApi{Apis: apis, Id: cli.id}, nil)
}

//多播
func (cli *GateCli) MultCast(msg interface{}, recvs ...uint64) error {
	return writeMsg(cli.conn, cli.codec, OP_MULTCAST, cli.id, msg, recvs)
}

//广播
func (cli *GateCli) BroadCast(msg interface{}) error {
	return writeMsg(cli.conn, cli.codec, OP_BROADCAST, cli.id, msg, nil)
}

//控制
func (cli *GateCli) Control(cmd uint32, args ...interface{}) error {
	return writeMsg(cli.conn, cli.codec, cmd, cli.id,
		&MsgControl{Cmd: cmd, Args: args}, nil)
}
