package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createKialiManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"alternateReportStorage": map[string]any{
				"storageClassName": generators.DebianStorageClass,
			},
			"trivy": map[string]any{
				"storageClassName": generators.DebianStorageClass,
				"resources": map[string]any{
					"limits": map[string]any{
						"memory": "1500M",
					},
				},
				"debug": true,
			},
			"operator": map[string]any{
				"scanJobsConcurrentLimit": 2,
				"logDevMode":              false,
			},
		},
		nil,
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
