package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createRepoManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, false),
	}

	repo := utils.GetGenericRepoManifest(istio.RepoName, generatorMeta.Helm)

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
		}),
	}
	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo})
}
