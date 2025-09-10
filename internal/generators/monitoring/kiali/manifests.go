package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createKialiManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, false),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"cr": map[string]any{
				"create":    true,
				"namespace": istio.Namespace,
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
