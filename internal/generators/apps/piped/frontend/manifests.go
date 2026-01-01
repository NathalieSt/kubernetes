package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createPipedManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"frontend": map[string]any{
				"service": map[string]any{
					"main": map[string]any{
						"enabled": "false",
					},
				},
				"env": map[string]any{
					"BACKEND_HOSTNAME": "piped-backend.cloud.nathalie-stiefsohn.eu",
				},
				"podSecurityContext": map[string]any{
					"runAsUser": 0,
				},
				"config": map[string]any{
					"PROXY_PART": "https://piped-ytproxy.cloud.nathalie-stiefsohn.eu",
				},
				"additionalContainers": map[string]any{
					"netbird-agent": core.Container{
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
			"backend": map[string]any{
				"enabled": false,
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
				"service": map[string]any{
					"backend": map[string]any{
						"enabled": "false",
					},
				},
			},
			"ytproxy": map[string]any{
				"enabled": false,
				"env":     map[string]any{},
			},
			"ingress": map[string]any{
				"main": map[string]any{
					"enabled": true,
				},
				"backend": map[string]any{
					"enabled": true,
				},
				"ytproxy": map[string]any{
					"enabled": true,
				},
			},
			"controller": map[string]any{
				"enabled": true,
			},
		},
		nil,
	)

	service := utils.ManifestConfig{
		Filename: "service.yaml",
		Manifests: []any{
			core.NewService(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Helm.Version,
					},
				}, core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/instance": "piped",
						"app.kubernetes.io/name":     "piped-frontend",
					},
					Ports: []core.ServicePort{
						{
							Name:       fmt.Sprintf("http-%v", generatorMeta.Name),
							Port:       generatorMeta.Port,
							TargetPort: generatorMeta.Port,
						},
					},
				},
			),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				repo.Filename,
				chart.Filename,
				release.Filename,
				service.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, repo, chart, release, service})
}
