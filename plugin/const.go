package plugin

import (
	"github.com/hanjingo/gate/com"
	ctl "github.com/hanjingo/gate/plugin/controller"
	flt "github.com/hanjingo/gate/plugin/filter"
	pm "github.com/hanjingo/gate/plugin/perm"
	st "github.com/hanjingo/gate/plugin/stream"
)

var Plugins map[string]com.PluginI

func Init() {
	Plugins = make(map[string]com.PluginI)

	//注册插件
	//过滤器
	filtv1 := flt.NewFilterV1()
	Plugins[filtv1.Name()] = filtv1

	//鉴权器
	permv1 := pm.NewPermV1()
	Plugins[permv1.Name()] = permv1

	//控流器
	streamv1 := st.NewStreamV1()
	Plugins[streamv1.Name()] = streamv1

	//控制器
	controllerv1 := ctl.NewControllerV1()
	Plugins[controllerv1.Name()] = controllerv1
}
