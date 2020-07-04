package main

import (
	"log"

	"github.com/seamounts/kubeapi/pkg/cli"
	pluginv1 "github.com/seamounts/kubeapi/pkg/plugin/v1"
)

func main() {
	c, err := cli.New(
		cli.WithPlugins(
			&pluginv1.Plugin{},
		),
		cli.WithDefaultPlugin(
			&pluginv1.Plugin{},
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
