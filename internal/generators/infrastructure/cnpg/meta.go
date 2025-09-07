package main

import (
	"kubernetes/pkg/schema/generator"
)

var CNPG = generator.GeneratorMeta{
	Name:          "cnpg",
	Namespace:     "cnpg-system",
	GeneratorType: generator.Infrastructure,
	Helm: generator.Helm{
		Chart:   "cloudnative-pg",
		Url:     "https://cloudnative-pg.github.io/charts",
		Version: "0.26.0",
	},
	DependsOnGenerators: []string{},
}
