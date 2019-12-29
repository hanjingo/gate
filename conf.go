package gate

type GateConfig struct {
	Id           uint8
	ConfServAddr string
	UserName     string
	PassWord     string
}

func NewGateConfig() *GateConfig {
	back := &GateConfig{}
	return back
}
