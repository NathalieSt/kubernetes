package main

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
)

var Glance = generator.GeneratorMeta{
	Name:          "glance",
	Namespace:     "glance",
	GeneratorType: generator.App,
	ClusterUrl:    "glance.glance.svc.cluster.local",
	Port:          9000,
	Docker: generator.Docker{
		Registry: "glanceapp/glance",
		//FIXME: set to nil, later fetch in generator from version.json
		Version: "v0.8.4",
	},
	Caddy: generator.Caddy{
		DNSName: "glance.cluster",
	},
	KedaScaling: keda.ScaledObjectTriggerMeta{
		Timezone:        "Europe/Vienna",
		Start:           "0 7 * * *",
		End:             "0 23 * * *",
		DesiredReplicas: "1",
	},
	DependsOnGenerators: []string{},
}
