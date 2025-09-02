package apps

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
)

var Mealie = generator.GeneratorMeta{
	Name:       "mealie",
	Namespace:  "mealie",
	EntityType: generator.App,
	ClusterUrl: "mealie.mealie.svc.cluster.local",
	Port:       9000,
	Docker: generator.Docker{
		Registry: "ghcr.io/mealie-recipes/mealie",
		//FIXME: set to nil, later fetch in generator from version.json
		Version: "v3.0.2",
	},
	KedaScaling: keda.ScaledObjectTriggerMeta{
		Timezone:        "Europe/Vienna",
		Start:           "0 9 * * *",
		End:             "0 21 * * *",
		DesiredReplicas: "1",
	},
	DependsOnGenerators: []string{},
}

var Jellyfin = generator.GeneratorMeta{
	Name:       "jellyfin",
	Namespace:  "jellyfin",
	EntityType: generator.App,
	ClusterUrl: "jellyfin.jellyfin.svc.cluster.local",
	Port:       9000,
	Helm: generator.Helm{
		Url:     "https://jellyfin.github.io/jellyfin-helm",
		Chart:   "jellyfin",
		Version: "2.3.0",
	},
	KedaScaling: keda.ScaledObjectTriggerMeta{
		Timezone:        "Europe/Vienna",
		Start:           "0 9 * * *",
		End:             "0 21 * * *",
		DesiredReplicas: "1",
	},
	DependsOnGenerators: []string{},
}

var Forgejo = generator.GeneratorMeta{
	Name:       "forgejo",
	Namespace:  "forgejo",
	EntityType: generator.App,
	ClusterUrl: "forgejo.forgejo.svc.cluster.local",
	Port:       9000,
	Helm: generator.Helm{
		Url:     "oci://code.forgejo.org/forgejo-helm/forgejo",
		Version: "14.0.0",
	},
	KedaScaling: keda.ScaledObjectTriggerMeta{
		Timezone:        "Europe/Vienna",
		Start:           "0 9 * * *",
		End:             "0 23 * * *",
		DesiredReplicas: "1",
	},
	DependsOnGenerators: []string{},
}
