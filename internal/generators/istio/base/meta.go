package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/pkg/schema/generator"
)

var Base = generator.GeneratorMeta{
	Name:          "base",
	Namespace:     istio.Namespace,
	GeneratorType: generator.Istio,
	Helm: &generator.Helm{
		Chart: "base",
	},
	DependsOnGenerators: []string{},
}
