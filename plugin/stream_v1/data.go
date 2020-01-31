package stream_v1

import (
	"time"

	"github.com/hanjingo/container"
)

func init() {
	agents = make(map[interface{}]*agentInfo)
}

var agents map[interface{}]*agentInfo

type agentInfo struct {
	id               interface{}     //id
	cache            *container.Pool //缓存
	life             time.Duration   //已使用的时长
	totalUsedStream  int64           //总共使用的流量(kb)
	usedStreamInDur  int64           //当前时间内已用流量(字节)
	totalStreamInDur int64           //当前时间内总可用流量
}
