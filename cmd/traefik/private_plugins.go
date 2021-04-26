package main

import (
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/plugins"
	"github.com/traefik/traefik/v2/pkg/server/middleware"
)

const privatePluginGoPath = "./private-plugins-storage/"

func createPluginBuilder(staticConfiguration *static.Configuration) (middleware.PluginsBuilder, error) {
	privatePlugins := make(map[string]plugins.Descriptor)
	if hasPlugins(staticConfiguration) {
		for k, v := range staticConfiguration.Experimental.Plugins {
			if v.Version != "private" {
				continue
			}

			privatePlugins[k] = v
			delete(staticConfiguration.Experimental.Plugins, k)
		}
	}

	pilotBuilder, err := createPilotPluginBuilder(staticConfiguration)
	if err != nil {
		return nil, err
	}

	if err := pilotBuilder.LoadPrivatePlugins(privatePluginGoPath, privatePlugins); err != nil {
		return nil, err
	}

	return pilotBuilder, nil
}
