package main

import (
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createNetbirdOperatorManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm, map[string]any{
		"managementURL": "https://netbird.nathalie-stiefsohn.eu",
		"netbirdAPI": map[string]any{
			"keyFromSecret": map[string]any{
				"name": shared.NetbirdAPIKeySecretName,
				"key":  "NB_API_KEY",
			},
		},
	}, nil)

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

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, namespace, repo, chart, release, release})
}
