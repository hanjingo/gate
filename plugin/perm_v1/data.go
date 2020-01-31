package perm_v1

func init() {
	agents = make(map[interface{}]*agentInfo)
}

var agents map[interface{}]*agentInfo

type agentInfo struct {
	id   interface{} //id
	perm uint32      //权限
}
