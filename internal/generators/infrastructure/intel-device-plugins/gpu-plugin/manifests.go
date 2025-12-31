package main

import (
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createVaultManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		nil,
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
