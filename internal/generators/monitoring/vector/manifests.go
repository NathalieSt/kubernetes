package main

import (
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createVectorManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"role": "Agent",
			"service": map[string]any{
				"enabled": false,
			},
			"customConfig": map[string]any{
				"data_dir": "/var/lib/vector",
				"api": map[string]any{
					"enabled": false,
				},
				"sources": map[string]any{
					"k8s_in": map[string]any{
						"type": "kubernetes_logs",
					},
				},
				"sinks": map[string]any{
					"es_cluster": map[string]any{
						"inputs": []string{
							"k8s_in",
						},
						"type": "elasticsearch",
						"endpoints": []string{
							"http://elasticsearch-es-internal-http.elastic-stack.svc.cluster.local:9200",
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
