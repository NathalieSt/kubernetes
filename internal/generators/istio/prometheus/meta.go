package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/pkg/schema/generator"
)

var Prometheus = generator.GeneratorMeta{
	Name:          "prometheus",
	Namespace:     istio.Namespace,
	GeneratorType: generator.Istio,
	Port:          9090,
	Docker: generator.Docker{
		Registry: "prom/prometheus",
		Version:  "v3.5.0",
	},
	DependsOnGenerators: []string{},
}
