package gate

//端点信息
type AgentInfo struct {
	CurrAgentNum     int //当前在线端点数
	HightestAgentNum int //同时在线最高端点数
}

//服务器信息
type ServerInfo struct {
	TcpNum int //tcp服务器数量
	WsNum  int //ws服务器数量
}

//本机信息
type MechineInfo struct {
	CpuPer float32 //cpu占用
	MemPer float32 //内存占用
	Volum  float32 //磁盘占用
}

//网关信息
type GateInfo struct {
	Agent   *AgentInfo
	Serv    *ServerInfo
	Mechine *MechineInfo
}

func newGateInfo() *GateInfo {
	back := &GateInfo{
		Agent:   &AgentInfo{},
		Serv:    &ServerInfo{},
		Mechine: &MechineInfo{},
	}
	return back
}
