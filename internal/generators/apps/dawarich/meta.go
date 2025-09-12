package main

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
)

var Dawarich = generator.GeneratorMeta{
	Name:          "dawarich",
	Namespace:     "dawarich",
	GeneratorType: generator.App,
	ClusterUrl:    "dawarich.dawarich.svc.cluster.local",
	Port:          3000,
	Docker: &generator.Docker{
		Registry: "freikin/dawarich",
		//FIXME: set to nil, later fetch in generator from version.json
		Version: "0.30.10",
	},
	Caddy: &generator.Caddy{
		DNSName: "dawarich.cluster",
	},
	KedaScaling: &keda.ScaledObjectTriggerMeta{
		Timezone:        "Europe/Vienna",
		Start:           "0 8 * * *",
		End:             "0 22 * * *",
		DesiredReplicas: "1",
	},
	DependsOnGenerators: []string{
		"redis",
		"postgres",
	},
}
