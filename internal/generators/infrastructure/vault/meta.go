package main

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
)

var Vault = generator.GeneratorMeta{
	Name:          "vault",
	Namespace:     "vault",
	GeneratorType: generator.Infrastructure,
	ClusterUrl:    "vault.vault.svc.cluster.local",
	Port:          9000,
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
