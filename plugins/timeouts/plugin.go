package timeouts

import (
	"time"

	"github.com/kennethklee/xpb"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func init() {
	xpb.Register(&Plugin{})
}

type Plugin struct{}

func (p *Plugin) Info() xpb.PluginInfo {
	return xpb.PluginInfo{
		Name:        "timeout",
		Version:     "latest",
		Description: "Timeouts DB queries after a certain amount of time",
	}
}

func (p *Plugin) OnPreload() error {
	return nil
}

func (p *Plugin) OnLoad(app core.App) error {
	pb := app.(*pocketbase.PocketBase)

	var queryTimeout int
	pb.RootCmd.PersistentFlags().IntVar(
		&queryTimeout,
		"queryTimeout",
		30,
		"the default SELECT queries timeout in seconds",
	)

	pb.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		app.Dao().ModelQueryTimeout = time.Duration(queryTimeout) * time.Second
		return nil
	})
	return nil
}
