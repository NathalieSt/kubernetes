package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
)

func createPipedManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"frontend": map[string]any{
				"enabled": false,
			},
			"backend": map[string]any{
				"config": map[string]any{
					"API_URL":      "https://piped-backend.cloud.nathalie-stiefsohn.eu",
					"PROXY_PART":   "https://piped-ytproxy.cloud.nathalie-stiefsohn.eu",
					"FRONTEND_URL": "https://piped.cloud.nathalie-stiefsohn.eu",
					"database": map[string]any{
						"connection_url": "jdbc:postgresql://postgres-rw.postgres.svc.cluster.local:5432/piped",
						"driver_class":   "org.postgresql.Driver",
						"dialect":        "org.hibernate.dialect.PostgreSQLDialect",
					},
				},
				"additionalContainers": []core.Container{
					{
						Name:  "netbird-agent",
						Image: "netbirdio/netbird:latest",
						Env: []core.Env{
							{
								Name: "NB_SETUP_KEY",
								ValueFrom: core.ValueFrom{
									SecretKeyRef: core.SecretKeyRef{
										Name: generators.NetbirdSecretName,
										Key:  "setup-key",
									},
								},
							},
							{
								Name:  "NB_MANAGEMENT_URL",
								Value: "https://netbird.nathalie-stiefsohn.eu",
							},
						},
						Resources: core.Resources{
							Requests: map[string]string{
								"cpu":    "50m",
								"memory": "64Mi",
							},
							Limits: map[string]string{
								"cpu":    "100m",
								"memory": "128Mi",
							},
						},
						SecurityContext: core.ContainerSecurityContext{
							Privileged: true,
						},
					},
				},
			},
			"ytproxy": map[string]any{
				"enabled": false,
			},
			"ingress": map[string]any{
				"main": map[string]any{
					"enabled": false,
				},
				"backend": map[string]any{
					"enabled": false,
				},
				"ytproxy": map[string]any{
					"enabled": false,
				},
			},
		},
		[]helm.ReleaseValuesFrom{
			{
				Kind:       helm.Secret,
				Name:       generators.PostgresCredsSecret,
				ValuesKey:  "username",
				TargetPath: "backend.config.database.username",
				Optional:   false,
			},
			{
				Kind:       helm.Secret,
				Name:       generators.PostgresCredsSecret,
				ValuesKey:  "password",
				TargetPath: "backend.config.database.password",
				Optional:   false,
			},
		},
	)

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				namespace.Filename,
				repo.Filename,
				chart.Filename,
				release.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release})
}
