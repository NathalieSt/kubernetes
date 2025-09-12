package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/pkg/schema/generator"
)

var Networking = generator.GeneratorMeta{
	Name:          "istio-networking",
	Namespace:     istio.Namespace,
	GeneratorType: generator.Istio,
	Helm: &generator.Helm{
		Url:     "oci://code.forgejo.org/forgejo-helm/forgejo",
		Version: "14.0.0",
	},
	DependsOnGenerators: []string{},
}
