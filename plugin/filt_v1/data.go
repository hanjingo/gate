package filt_v1

func init() {
	agents = make(map[interface{}]*agentInfo)
}

var agents map[interface{}]*agentInfo

//端点
type agentInfo struct {
	id    interface{} //id
	bFilt bool        //是否被过滤
}
