package main

import (
	"kubernetes/pkg/schema/generator"
)

var Ingress = generator.GeneratorMeta{
	Name:          "istio-ingress",
	Namespace:     "istio-ingress",
	GeneratorType: generator.Istio,
	Helm: generator.Helm{
		Chart:   "gateway",
		Version: "1.27.0",
	},
	DependsOnGenerators: []string{},
}
