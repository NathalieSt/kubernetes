package main

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
)

var Jellyfin = generator.GeneratorMeta{
	Name:          "jellyfin",
	Namespace:     "jellyfin",
	GeneratorType: generator.App,
	ClusterUrl:    "jellyfin.jellyfin.svc.cluster.local",
	Port:          8096,
	Helm: generator.Helm{
		Url:     "https://jellyfin.github.io/jellyfin-helm",
		Chart:   "jellyfin",
		Version: "2.3.0",
	},
	Caddy: generator.Caddy{
		DNSName:                    "jellyfin.cluster",
		WebsocketSupportIsRequired: true,
	},
	KedaScaling: keda.ScaledObjectTriggerMeta{
		Timezone:        "Europe/Vienna",
		Start:           "0 9 * * *",
		End:             "0 23 * * *",
		DesiredReplicas: "1",
	},
	DependsOnGenerators: []string{},
}
