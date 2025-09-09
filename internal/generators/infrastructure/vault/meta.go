package main

import (
	"kubernetes/pkg/schema/generator"
)

var Vault = generator.GeneratorMeta{
	Name:          "vault",
	Namespace:     "vault",
	GeneratorType: generator.Infrastructure,
	ClusterUrl:    "vault.vault.svc.cluster.local",
	Port:          9000,
	Helm: generator.Helm{
		Chart:   "vault",
		Url:     "https://helm.releases.hashicorp.com",
		Version: "0.30.1",
	},
	DependsOnGenerators: []string{},
}
