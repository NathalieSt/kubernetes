package main

import (
	"kubernetes/pkg/schema/generator"
)

var Keda = generator.GeneratorMeta{
	Name:          "keda",
	Namespace:     "keda",
	GeneratorType: generator.Infrastructure,
	Helm: generator.Helm{
		Chart:   "keda",
		Url:     "https://kedacore.github.io/charts",
		Version: "2.17.2",
	},
	DependsOnGenerators: []string{},
}
