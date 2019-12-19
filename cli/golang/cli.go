package golang

import (
	"errors"
	"fmt"
	"time"

	"github.com/hanjingo/gate/com"
	ctlv1 "github.com/hanjingo/gate/plugin/control_v1"
	"github.com/hanjingo/util"

	ws "github.com/gorilla/websocket"
	"github.com/hanjingo/network"
	"github.com/hanjingo/protocol"
	pv4 "github.com/hanjingo/protocol/v4"
)

type CallFunc func(interface{})

func writeMsg(conn network.ConnI, codec protocol.CodecI, opcode uint32, sender uint64, content interface{}, recvs []uint64) error {
	p := &pv4.Messager{
		OpCode:    opcode,
		Sender:    sender,
		Receiver:  recvs,
		TimeStamp: util.TimeDurToMilliSecond(time.Now().Sub(com.START_TIME)),
		Content:   content,
	}
	data, err := codec.Format(p)
	if err != nil {
		return err
	}
	return conn.WriteMsg(data)
}

func readMsg(conn network.ConnI, codec protocol.CodecI, content interface{}) error {
	data, err := conn.ReadMsg()
	if err != nil {
		return err
	}
	return codec.UnFormat(data, content)
}

//网关客户端
type GateCli struct {
	id                uint64
	conn              network.ConnI
	callMap           map[uint32]CallFunc
	codec             protocol.CodecI
	connCloseCallback func(network.ConnI)
}

func NewGateCli(f func(network.ConnI)) *GateCli {
	back := &GateCli{
		callMap:           make(map[uint32]CallFunc),
		codec:             pv4.NewCodec(),
		connCloseCallback: f,
	}
	return back
}

//注册回调
func (cli *GateCli) RegHandler(opcdoe uint32, f CallFunc) error {
	if _, ok := cli.callMap[opcdoe]; ok {
		return errors.New("已经注册,无需再次注册")
	}
	cli.callMap[opcdoe] = f
	return nil
}

//拨号
func (cli *GateCli) Dial(dialType string, addr string, token string, conf *network.ConnConfig) error {
	var conn network.ConnI
	var err error
	c, _, err := ws.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	conn, err = network.NewWsConn(conf, c, cli.onConnClose, nil)
	if err != nil {
		return err
	}
	//发请求
	if err := writeMsg(conn, cli.codec, OP_NEW_AGENT, 0, &ctlv1.MsgNewAgentReq{Token: token}, nil); err != nil {
		return err
	}
	//收回复
	rsp := &ctlv1.MsgNewAgentRsp{}
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

func (cli *GateCli) onConnClose(c network.ConnI) {
	if cli.connCloseCallback != nil {
		cli.connCloseCallback(c)
	}
}

//路由
func (cli *GateCli) Route(msg interface{}, recvs ...uint64) error {
	return writeMsg(cli.conn, cli.codec, OP_ROUTE, cli.id, msg, recvs)
}

//ping
func (cli *GateCli) Ping() error {
	return writeMsg(cli.conn, cli.codec, OP_PING, cli.id,
		&ctlv1.MsgPing{}, nil)
}

//请求
func (cli *GateCli) Req(opCode uint32, msg interface{}, recv ...uint64) error {
	return writeMsg(cli.conn, cli.codec, MASK_CLI&opCode, cli.id, msg, recv)
}

//订阅
func (cli *GateCli) Sub(topics ...string) error {
	return writeMsg(cli.conn, cli.codec, OP_SUB, cli.id,
		&ctlv1.MsgSub{Topics: topics}, nil)
}

//发布
func (cli *GateCli) Pub() {

}

//取消订阅
func (cli *GateCli) UnSub(topics ...string) error {
	return writeMsg(cli.conn, cli.codec, OP_UNSUB, cli.id,
		&ctlv1.MsgUnSub{Topics: topics}, nil)
}

//注册服务
func (cli *GateCli) RegApi(apis ...uint32) error {
	return writeMsg(cli.conn, cli.codec, OP_REG_SERVER, cli.id,
		&ctlv1.MsgRegApi{Apis: apis, Id: cli.id}, nil)
}

//取消注册
func (cli *GateCli) UnRegApi(apis ...uint32) error {
	return writeMsg(cli.conn, cli.codec, OP_UNREG_SERVER, cli.id,
		&ctlv1.MsgUnRegApi{Apis: apis, Id: cli.id}, nil)
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
		&ctlv1.MsgControl{Cmd: cmd, Args: args}, nil)
}
