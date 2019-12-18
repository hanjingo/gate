package plugin

import (
	"github.com/hanjingo/gate/com"
	ctlv1 "github.com/hanjingo/gate/plugin/control_v1"
	filtv1 "github.com/hanjingo/gate/plugin/filt_v1"
	permv1 "github.com/hanjingo/gate/plugin/perm_v1"
	streamv1 "github.com/hanjingo/gate/plugin/stream_v1"
)

var pm map[string]com.PluginI

func init() {
	if pm == nil {
		pm = make(map[string]com.PluginI)
	}
	//控制器 版本:1
	reg(ctlv1.NewControllerV1())
	//过滤器 版本:1
	reg(filtv1.NewFilt())
	//鉴权器 版本:1
	reg(permv1.NewPerm())
	//控流器 版本:1
	reg(streamv1.NewStream())
}

func GetPluginMap() map[string]com.PluginI {
	return pm
}

func reg(p com.PluginI) {
	pm[p.Name()] = p
}
