package xpb

import (
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

type Plugin interface {
	/**
	 * Preload is called before the app is setup.
	 * This is a good place to load configurations.
	 */
	OnPreload() error

	/**
	 * Load is called after the app is setup.
	 * This is a good place to register commands
	 * and hooks.
	 */
	OnLoad(app core.App) error

	/**
	 * Get plugin info
	 */
	Info() PluginInfo
}

// For display purposes only
type PluginInfo struct {
	Name        string
	Version     string
	Description string
}

var plugins = []Plugin{}

func Register(plugin Plugin) {
	plugins = append(plugins, plugin)
}

func FireOnPreload() (err error) {
	for _, plugin := range plugins {
		if pluginErr := plugin.OnPreload(); pluginErr != nil {
			err = errors.Join(err, pluginErr)
		}
	}
	return
}

func FireOnLoad(app core.App) (err error) {
	for _, plugin := range plugins {
		if pluginErr := plugin.OnLoad(app); pluginErr != nil {
			err = errors.Join(err, pluginErr)
		}
	}
	return
}

func GetPlugins() []Plugin {
	return plugins
}
