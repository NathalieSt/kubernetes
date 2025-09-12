package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/pkg/schema/generator"
)

var Istiod = generator.GeneratorMeta{
	Name:          "istiod",
	Namespace:     istio.Namespace,
	GeneratorType: generator.Istio,
	Helm: &generator.Helm{
		Chart:   "istiod",
		Version: "1.27.0",
	},
	DependsOnGenerators: []string{},
}
