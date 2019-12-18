package golang

import (
	"fmt"
	"errors"
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

func newProto() *pv4.Messager {
	return &pv4.Messager{
		OpCode:    0,
		Receiver:  []uint64{},
		Sender:    0,
		TimeStamp: util.TimeDurToMilliSecond(time.Now().Sub(com.START_TIME)),
		Content:   nil,
	}
}

//网关客户端
type GateCli struct {
	id      uint64
	conn    network.ConnI
	callMap map[uint32]CallFunc
	codec   protocol.CodecI
}

func NewGateCli() *GateCli {
	back := &GateCli{
		callMap: make(map[uint32]CallFunc),
		codec:   pv4.NewCodec(),
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
	conn.Run()
	//发请求
	p := newProto()
	p.OpCode = OP_NEW_AGENT
	p.Content = &ctlv1.MsgNewAgentReq{Token: token}
	data, err := cli.codec.Format(p)
	if err != nil {
		return err
	}
	if err := conn.WriteMsg(data); err != nil {
		return err
	}
	fmt.Println("我在这1")
	//收回复
	data1, err := conn.ReadMsg()
	if err != nil {
		return err
	}
	p1 := &pv4.Messager{Content: &ctlv1.MsgNewAgentRsp{}}
	if err := cli.codec.UnFormat(data1, p1); err != nil {
		return err
	}
	rsp := p1.Content.(*ctlv1.MsgNewAgentRsp)
	if !rsp.Result {
		return errors.New("网关拒绝建立连接")
	}
	cli.conn = conn
	cli.id = rsp.Id
	fmt.Println("与网关:", addr, "连接建立")
	return nil
}

func (cli *GateCli) onConnClose(c network.ConnI) {
	//todo
}

//路由
func (cli *GateCli) Route(msg interface{}, recvs ...uint64) error {
	p := newProto()
	p.OpCode = OP_ROUTE
	p.Receiver = recvs
	p.Sender = cli.id
	p.Content = msg
	data, err := cli.codec.Format(p)
	if err != nil {
		return err
	}
	return cli.conn.WriteMsg(data)
}

//ping
func (cli *GateCli) Ping(msg interface{}) error {
	p := newProto()
	p.OpCode = OP_PING
	p.Sender = cli.id
	p.Content = msg
	data, err := cli.codec.Format(p)
	if err != nil {
		return err
	}
	return cli.conn.WriteMsg(data)
}

//请求
func (cli *GateCli) Req(opCode uint32, msg interface{}, recv ...uint64) error {
	return nil
}

//订阅
func (cli *GateCli) Sub() {

}

//取消订阅
func (cli *GateCli) UnSub() {

}

//注册服务
func (cli *GateCli) RegApi() {

}

//取消注册
func (cli *GateCli) UnRegApi() {

}

//多播
func (cli *GateCli) MultCast() {

}

//广播
func (cli *GateCli) BroadCast() {

}

//发布
func (cli *GateCli) Pub() {

}

//控制
func (cli *GateCli) Control() {

}
