package gate

import (
	"errors"

	"github.com/hanjingo/gate/com"
	msg "github.com/hanjingo/msg"
)

type codec struct {
	versions map[uint8]msg.MsgI //key:编解码器id value:编解码器生成函数
}

func newCodec() *codec {
	back := &codec{
		versions: make(map[uint8]msg.MsgI),
	}
	back.versions[com.CODEC_ID1] = msg.GetMsger()[msg.VERSION1]
	back.versions[com.CODEC_ID2] = msg.GetMsger()[msg.VERSION2]
	back.versions[com.CODEC_ID3] = msg.GetMsger()[msg.VERSION3]
	back.versions[com.CODEC_ID4] = msg.GetMsger()[msg.VERSION4]
	return back
}

func (c *codec) getMsger(codecId uint8) msg.MsgI {
	return c.versions[codecId]
}

func (c *codec) UnFormat(data []byte) (*com.Msg, error) {
	if data == nil || len(data) < 1 {
		return nil, errors.New("消息过短")
	}
	msger := c.getMsger(data[0])
	if msger == nil {
		return nil, errors.New("未知编码格式的消息")
	}
	back := com.NewMsg()
	_, err := msger.UnFormat(data[1:], back)
	return back, err
}

func (c *codec) Format(content interface{}, codecId uint8) ([]byte, error) {
	if msger := c.getMsger(codecId); msger != nil {
		return msger.Format(content)
	}
	return nil, errors.New("未知编码格式")
}
