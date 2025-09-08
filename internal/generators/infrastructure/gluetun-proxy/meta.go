package main

import (
	"kubernetes/pkg/schema/generator"
)

var GluetunProxy = generator.GeneratorMeta{
	Name:          "gluetun-proxy",
	Namespace:     "gluetun-proxy",
	GeneratorType: generator.Infrastructure,
	ClusterUrl:    "gluetun-proxy.gluetun-proxy.svc.cluster.local",
	Port:          8888,
	Docker: generator.Docker{
		Registry: "qmcgaw/gluetun",
		Version:  "v3.40",
	},
	DependsOnGenerators: []string{},
}
