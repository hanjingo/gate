package plugin

import (
	ps "github.com/hanjingo/golib/plugin_system"

	ctlv1 "github.com/hanjingo/gate/plugin/control_v1"
	filtv1 "github.com/hanjingo/gate/plugin/filt_v1"
	permv1 "github.com/hanjingo/gate/plugin/perm_v1"
	streamv1 "github.com/hanjingo/gate/plugin/stream_v1"
)

var plugins map[string]func() ps.PluginI

func GetPlugins() map[string]func() ps.PluginI {
	if plugins == nil {
		plugins = make(map[string]func() ps.PluginI)
	}
	return plugins
}

func init() {
	GetPlugins()[ctlv1.NAME] = ctlv1.New
	GetPlugins()[filtv1.NAME] = filtv1.New
	GetPlugins()[permv1.NAME] = permv1.New
	GetPlugins()[streamv1.NAME] = streamv1.New
}
