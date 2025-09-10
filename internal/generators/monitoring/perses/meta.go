package main

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
)

var Perses = generator.GeneratorMeta{
	Name:          "perses",
	Namespace:     "perses",
	GeneratorType: generator.Monitoring,
	ClusterUrl:    "perses.perses.svc.cluster.local",
	Port:          8080,
	Docker: generator.Docker{
		Registry: "ghcr.io/mealie-recipes/mealie",
		//FIXME: set to nil, later fetch in generator from version.json
		Version: "v3.0.2",
	},
	Caddy: generator.Caddy{
		DNSName: "perses.cluster",
	},
	KedaScaling: keda.ScaledObjectTriggerMeta{
		Timezone:        "Europe/Vienna",
		Start:           "0 8 * * *",
		End:             "0 22 * * *",
		DesiredReplicas: "1",
	},
	DependsOnGenerators: []string{},
}
