package control_v1

import (
	"time"

	"github.com/hanjingo/gate/com"
	"github.com/hanjingo/golib/container"
)

func init() {
	agents = make(map[interface{}]*agentInfo)
}

var agents map[interface{}]*agentInfo

type agentInfo struct {
	agent     com.AgentI
	id        interface{}    //id
	startTm   time.Time      //开始时间
	subTopics *container.Set //订阅的主题
	apis      *container.Set //可提供的api服务
}
