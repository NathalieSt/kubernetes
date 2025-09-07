package main

import (
	"kubernetes/pkg/schema/generator"
)

var Caddy = generator.GeneratorMeta{
	Name:          "caddy",
	Namespace:     "caddy",
	GeneratorType: generator.Infrastructure,
	ClusterUrl:    "caddy.caddy.svc.cluster.local",
	Port:          80,
	Docker: generator.Docker{
		Registry: "caddy",
		Version:  "2.10.0-alpine",
	},
	DependsOnGenerators: []string{},
}
