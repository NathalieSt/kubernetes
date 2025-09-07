package main

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
)

var SearXNG = generator.GeneratorMeta{
	Name:          "searxng",
	Namespace:     "searxng",
	GeneratorType: generator.App,
	ClusterUrl:    "searxng.searxng.svc.cluster.local",
	Port:          8080,
	Docker: generator.Docker{
		Registry: "searxng/searxng",
		//FIXME: set to nil, later fetch in generator from version.json
		Version: "2025.8.3-2e62eb5",
	},
	Caddy: generator.Caddy{
		DNSName: "searxng.cluster",
	},
	KedaScaling: keda.ScaledObjectTriggerMeta{
		Timezone:        "Europe/Vienna",
		Start:           "0 7 * * *",
		End:             "0 23 * * *",
		DesiredReplicas: "1",
	},
	DependsOnGenerators: []string{},
}
