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
}
