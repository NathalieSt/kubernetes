package main

import (
	"kubernetes/internal/generators"
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
		"ingress": map[string]any{
			"enabled": true,
			"router": map[string]any{
				"enabled": false,
			},
		},
		"netbirdAPI": map[string]any{
			"keyFromSecret": map[string]any{
				"name": generators.NetbirdAPIKeySecretName,
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

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release})
}
