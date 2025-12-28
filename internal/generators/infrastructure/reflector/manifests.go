package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/vaultsecretsoperator"
	"kubernetes/pkg/schema/generator"
)

func createReflectorManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm, nil, nil)

	netbirdSecretConfig := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-netbird-static-secret", generatorMeta.Name),
		SecretName: generators.NetbirdSecretName,
		Path:       "netbird/setup-key",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "caddy",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "caddy",
		},
	}

	hetznerAPITokenSecret := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-hetzner-api-token-static-secret", generatorMeta.Name),
		SecretName: generators.HetznerAPITokenSecretName,
		Path:       "hetzner/api-token",
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

	forgejoPGSecretConfig := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-forgejo-pg-static-secret", generatorMeta.Name),
		SecretName: generators.ForgejoPGCredsSecret,
		Path:       "postgres-clusters/forgejo",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "forgejo-pg-cluster,forgejo",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "forgejo-pg-cluster,forgejo",
		},
	}

	matrixPGSecretConfig := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-matrix-pg-static-secret", generatorMeta.Name),
		SecretName: generators.MatrixPGCredsSecret,
		Path:       "postgres-clusters/matrix",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "matrix-pg-cluster,synapse,discord-bridge",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "matrix-pg-cluster,synapse,discord-bridge",
		},
	}

	mariaDBSecretConfig := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-mariadb-static-secret", generatorMeta.Name),
		SecretName: generators.MariaDBCredsSecret,
		Path:       "mariadb",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "mariadb,booklore",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "mariadb,booklore",
		},
	}

	synapseSecret := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-synapse-static-secret", generatorMeta.Name),
		SecretName: generators.SynapseSecretName,
		Path:       "matrix/synapse",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "synapse",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "synapse",
		},
	}

	discordBridgeSecret := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-discord-bridge-static-secret", generatorMeta.Name),
		SecretName: generators.DiscordBridgeSecretName,
		Path:       "matrix/discord-bridge",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "synapse,discord-bridge",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "synapse,discord-bridge",
		},
	}

	whatsappBridgeSecret := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-whatsapp-bridge-static-secret", generatorMeta.Name),
		SecretName: generators.WhatsappBridgeSecretName,
		Path:       "matrix/whatsapp-bridge",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "synapse,whatsapp-bridge",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "synapse,whatsapp-bridge",
		},
	}

	elasticAdminSecret := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-elastic-search-admin-static-secret", generatorMeta.Name),
		SecretName: generators.ElasticSearchAdminSecretName,
		Path:       "elastic-search/admin",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "elastic-stack",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "elastic-stack",
		},
		Type: vaultsecretsoperator.BasicAuth,
	}

	elasticVectorSecret := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-elastic-search-vector-static-secret", generatorMeta.Name),
		SecretName: generators.ElasticSearchVectorSecretName,
		Path:       "elastic-search/vector",
		SecretAnnotations: map[string]string{
			"reflector.v1.k8s.emberstack.com/reflection-allowed":            "true",
			"reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "elastic-stack,vector",
			"reflector.v1.k8s.emberstack.com/reflection-auto-enabled":       "true",
			"reflector.v1.k8s.emberstack.com/reflection-auto-namespaces":    "elastic-stack,vector",
		},
		Type: vaultsecretsoperator.BasicAuth,
	}

	vaultSecrets := utils.ManifestConfig{
		Filename: "vault-secrets.yaml",
		Manifests: utils.GenerateVaultAccessManifests(
			generatorMeta.Name,
			//FIXME: get this from VSO generator meta
			"vault-secrets-operator",
			[]utils.StaticSecretConfig{
				netbirdSecretConfig,
				hetznerAPITokenSecret,
				postgresSecretConfig,
				forgejoPGSecretConfig,
				matrixPGSecretConfig,
				mariaDBSecretConfig,
				synapseSecret,
				discordBridgeSecret,
				whatsappBridgeSecret,
				elasticAdminSecret,
				elasticVectorSecret,
			},
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
