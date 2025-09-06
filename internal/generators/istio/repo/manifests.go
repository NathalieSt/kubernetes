package main

import (
	"kubernetes/internal/generators/istio"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createRepoManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, false),
	}

	repo := utils.ManifestConfig{
		Filename: "repo.yaml",
		Manifests: []any{
			helm.NewRepo(meta.ObjectMeta{
				Name: istio.RepoName,
			}, helm.RepoSpec{
				RepoType: helm.Default,
				Url:      generatorMeta.Helm.Url,
				Interval: "24h",
			}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
		}),
	}
	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo})
}
