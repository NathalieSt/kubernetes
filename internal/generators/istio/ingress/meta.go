package main

import (
	"kubernetes/pkg/schema/generator"
)

var Ingress = generator.GeneratorMeta{
	Name:          "istio-ingress",
	Namespace:     "istio-ingress",
	GeneratorType: generator.Istio,
	Helm: generator.Helm{
		Url:     "https://istio-release.storage.googleapis.com/charts",
		Chart:   "gateway",
		Version: "1.27.0",
	},
	DependsOnGenerators: []string{},
}
