package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/pkg/schema/generator"
)

var security = generator.GeneratorMeta{
	Name:                "istio-security",
	Namespace:           istio.Namespace,
	GeneratorType:       generator.Istio,
	DependsOnGenerators: []string{},
}
