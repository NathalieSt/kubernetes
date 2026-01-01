package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
)

func createPipedManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"frontend": map[string]any{
				"podSecurityContext": map[string]any{
					"runAsUser": 0,
				},
				"env": map[string]any{
					"BACKEND_HOSTNAME": "piped-backend.cloud.nathalie-stiefsohn.eu",
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
			"backend": map[string]any{
				"enabled": false,
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
		nil,
	)

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				repo.Filename,
				chart.Filename,
				release.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, repo, chart, release})
}
