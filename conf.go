package gate

type GateConfig struct {
	Id           uint8
	ConfServAddr string
	UserName     string
	PassWord     string
	Plugins      map[string]*PluginConfig
}

type PluginConfig struct {
	Name string
}

func NewGateConfig() *GateConfig {
	back := &GateConfig{
		Plugins: make(map[string]*PluginConfig),
	}
	return back
}
