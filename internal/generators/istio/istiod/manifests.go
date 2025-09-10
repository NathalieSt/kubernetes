package main

import (
	"fmt"
	"kubernetes/internal/generators/istio"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/kustomize"
)

func createIstiodManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	chartName := fmt.Sprintf("%v-chart", generatorMeta.Name)
	chart := utils.GetGenericChartManifest(chartName, generatorMeta.Helm, istio.RepoName)

	release := utils.GetGenericReleaseManifest(generatorMeta.Name, chartName, nil, nil)

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: []any{
			kustomize.NewKustomization(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
				},
				[]string{
					chart.Filename,
					release.Filename,
				},
			),
		},
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, chart, release})
}
