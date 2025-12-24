package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"os"
)

func createElasticStackManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	config, err := os.ReadFile("./beats-config.yaml")
	if err != nil {
		fmt.Printf("Error while reading config.yaml")
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"elasticsearch": map[string]any{
				"spec": map[string]any{
					"nodeSets": map[string]any{
						"volumeClaimTemplates": []map[string]any{
							{
								"metadata": map[string]any{
									"name": "elasticsearch-data",
								},
								"spec": map[string]any{
									"accessModes": []string{
										"ReadWriteOnce",
									},
									"storageClassName": generators.NFSRemoteClass,
									"resources": map[string]any{
										"requests": map[string]any{
											"storage": "20Gi",
										},
									},
								},
							},
						},
					},
				},
			},
			"eck-beats": map[string]any{
				"enabled": true,
				"type":    "filebeat",
				"elasticsearchRef": map[string]any{
					"name": "elasticsearch",
				},
				"kibanaRef": map[string]any{
					"name": "elastic-stack-instance-eck-kibana",
				},
				"config": map[string]any{
					"filebeat.yml": string(config),
				},
				"daemonSet": map[string]any{
					"podTemplate": map[string]any{
						"spec": map[string]any{
							"securityContext": map[string]any{
								"runAsUser": 0,
							},
						},
					},
				},
			},
		},
		nil,
	)

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
			chart.Filename,
			release.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release})
}
