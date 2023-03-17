package migrations

import (
	"os"

	"github.com/kennethklee/xpb"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func init() {
	xpb.Register(&Plugin{})
}

type Plugin struct{}

func (p *Plugin) Info() xpb.PluginInfo {
	return xpb.PluginInfo{
		Name:        "migrations",
		Version:     "latest",
		Description: "Database migrations",
	}
}

func (p *Plugin) OnPreload() error {
	return nil
}

func (p *Plugin) OnLoad(app core.App) error {
	pb := app.(*pocketbase.PocketBase)

	var migrationsDir string
	pb.RootCmd.PersistentFlags().StringVar(
		&migrationsDir,
		"migrationsDir",
		"",
		"the directory with the user defined migrations",
	)

	var automigrate bool
	pb.RootCmd.PersistentFlags().BoolVar(
		&automigrate,
		"automigrate",
		true,
		"enable/disable auto migrations",
	)

	pb.RootCmd.ParseFlags(os.Args[1:])

	// migrate command (for go)
	migratecmd.MustRegister(app, pb.RootCmd, &migratecmd.Options{
		Automigrate: automigrate,
		Dir:         migrationsDir,
	})

	return nil
}
