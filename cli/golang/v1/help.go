package v1

import (
	"github.com/hanjingo/golib/network"
	"github.com/hanjingo/golib/protocol"
	pv4 "github.com/hanjingo/golib/protocol/v4"
)

func writeMsg(conn network.SessionI, codec protocol.CodecI, opcode interface{}, sender uint64, content interface{}, recvs []uint64) error {
	msg := &pv4.Messager{
		OpCode:   opcode.(uint32),
		Sender:   sender,
		Receiver: recvs,
		Content:  content,
	}
	data, err := codec.Format(msg)
	if err != nil {
		return err
	}
	return conn.WriteMsg(data)
}

func readMsg(conn network.SessionI, codec protocol.CodecI, content interface{}) error {
	data, err := conn.ReadMsg()
	if err != nil {
		return err
	}
	msg := &pv4.Messager{
		Content: content,
	}
	return codec.UnFormat(data, msg)
}
