package main

import (
	"fmt"
	"kubernetes/internal/generators/istio"
	"kubernetes/pkg/schema/generator"
)

var Kiali = generator.GeneratorMeta{
	Name:          "kiali",
	Namespace:     "kiali-operator",
	GeneratorType: generator.Monitoring,
	ClusterUrl:    fmt.Sprintf("kiali.%v.svc.cluster.local", istio.Namespace),
	Port:          20001,
	Caddy: &generator.Caddy{
		DNSName: "kiali.cluster",
	},
	Helm: &generator.Helm{
		Chart:   "kiali-operator",
		Url:     "https://kiali.org/helm-charts",
		Version: "2.14.0",
	},
	DependsOnGenerators: []string{},
}
