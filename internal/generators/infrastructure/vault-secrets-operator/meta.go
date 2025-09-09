package main

import (
	"kubernetes/pkg/schema/generator"
)

var VaultSecretsOperator = generator.GeneratorMeta{
	Name:          "vault-secrets-operator",
	Namespace:     "vault-secrets-operator",
	GeneratorType: generator.Infrastructure,
	Helm: generator.Helm{
		Chart:   "vault-secrets-operator",
		Url:     "https://helm.releases.hashicorp.com",
		Version: "0.10.0",
	},
	DependsOnGenerators: []string{},
}
