package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createReflectorManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm, nil, nil)

	netbirdSecretConfig := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-netbird-static-secret", generatorMeta.Name),
		SecretName: generators.NetbirdNecretName,
		Path:       "netbird/setup-key",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "caddy",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "caddy",
		},
	}

	postgresSecretConfig := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-postgres-static-secret", generatorMeta.Name),
		SecretName: generators.PostgresCredsSecret,
		Path:       "postgres",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "postgres,dawarich,mealie,forgejo,keycloak",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "postgres,dawarich,mealie,forgejo,keycloak",
		},
	}

	vaultSecrets := utils.ManifestConfig{
		Filename: "vault-secrets.yaml",
		Manifests: utils.GenerateVaultAccessManifests(
			generatorMeta.Name,
			//FIXME: get this from VSO generator meta
			"vault-secrets-operator",
			[]utils.StaticSecretConfig{netbirdSecretConfig, postgresSecretConfig},
		),
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
			chart.Filename,
			release.Filename,
			vaultSecrets.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, vaultSecrets})
}
