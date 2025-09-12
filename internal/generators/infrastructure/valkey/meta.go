package main

import (
	"kubernetes/pkg/schema/generator"
)

var Valkey = generator.GeneratorMeta{
	Name:          "valkey",
	Namespace:     "valkey",
	GeneratorType: generator.Infrastructure,
	ClusterUrl:    "valkey.valkey.svc.cluster.local",
	Port:          6379,
	Docker: &generator.Docker{
		Registry: "valkey/valkey",
		Version:  "8-alpine3.22",
	},
	DependsOnGenerators: []string{},
}
