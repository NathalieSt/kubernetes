package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/pkg/schema/generator"
)

var Repo = generator.GeneratorMeta{
	Name:          "istio-repo",
	Namespace:     istio.Namespace,
	GeneratorType: generator.Istio,
	Helm: generator.Helm{
		Url: "https://istio-release.storage.googleapis.com/charts",
	},
	DependsOnGenerators: []string{},
}
