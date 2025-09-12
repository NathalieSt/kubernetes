package main

import (
	"kubernetes/pkg/schema/generator"
)

var Reflector = generator.GeneratorMeta{
	Name:          "reflector",
	Namespace:     "reflector",
	GeneratorType: generator.Infrastructure,
	Helm: &generator.Helm{
		Chart:   "reflector",
		Url:     "https://emberstack.github.io/helm-charts",
		Version: "9.1.27",
	},
	DependsOnGenerators: []string{},
}
