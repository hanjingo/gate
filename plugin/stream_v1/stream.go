package stream_v1

import (
	"time"

	ps "github.com/hanjingo/golib/plugin"
	"github.com/hanjingo/golib/util"
)

const NAME = "StreamV1" //插件名字
const VERSION = "1.0.0" //插件版本

type StreamV1 struct {
	name        string
	info        *ps.PluginInfo       //插件消息
	streamSpeed int                  //流速(字节/s)
	scanDur     int                  //扫描间隔(s)
	reSendMsg   func(uint64, []byte) //重发消息函数
}

//streamSize:默认流速 	scanDur:扫描间隔	f:消息重发函数
func New() ps.PluginI {
	back := &StreamV1{
		name:        NAME,
		streamSpeed: int(util.MB * 64), //默认流速
		scanDur:     10,                //扫描间隔
		reSendMsg:   nil,
	}
	back.info = &ps.PluginInfo{
		Id:          back.name,
		Type:        ps.PLUGIN_TYPE_MEM,
		Version:     VERSION,
		Objs:        make(map[string]*ps.Object),
		CallBackMap: make(map[interface{}]interface{}),
	}
	back.reg()
	back.run()
	return back
}

func (s *StreamV1) Info() *ps.PluginInfo {
	return s.info
}

func (s *StreamV1) agents() map[interface{}]*agentInfo {
	return agents
}

//跑起来
func (s *StreamV1) run() {
	go func() {
		scanDur := time.Duration(s.scanDur) * time.Second
		tm := time.NewTimer(scanDur)
		for {
			select {
			case <-tm.C:
				for _, info := range s.agents() {
					info.usedStreamInDur = 0
					info.life += scanDur
					//把消息取出来发给插件总线
					for data := info.cache.Get(); data != nil; {
						if s.reSendMsg != nil {
							s.reSendMsg(info.id.(uint64), data.([]byte))
						}
					}
				}
				tm.Reset(scanDur)
			}
		}
	}()
}
