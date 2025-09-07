package main

import (
	"kubernetes/pkg/schema/generator"
)

var Redis = generator.GeneratorMeta{
	Name:          "redis",
	Namespace:     "redis",
	GeneratorType: generator.Infrastructure,
	ClusterUrl:    "redis.redis.svc.cluster.local",
	Port:          6379,
	Docker: generator.Docker{
		Registry: "redis",
		Version:  "8.2.1",
	},
	DependsOnGenerators: []string{},
}
