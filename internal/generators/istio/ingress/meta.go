package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/pkg/schema/generator"
)

var Ingress = generator.GeneratorMeta{
	Name:          "istio-ingress",
	Namespace:     istio.Namespace,
	GeneratorType: generator.Istio,
	Helm: generator.Helm{
		Chart:   "gateway",
		Version: "1.27.0",
	},
	DependsOnGenerators: []string{},
}
