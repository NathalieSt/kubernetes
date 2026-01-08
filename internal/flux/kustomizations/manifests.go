package main

import (
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/kustomization"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"sort"
	"sync"
)

func metaToFluxKustomizationManifest(generatorMeta generator.GeneratorMeta) kustomization.Kustomization {
	return kustomization.NewKustomization(meta.ObjectMeta{
		Name:      generatorMeta.Name,
		Namespace: "flux-system",
	}, *generatorMeta.Flux)
}

func createFluxKustomizationManifests(rootDir string) map[string][]byte {

	metas, err := utils.GetDiscoveredGeneratorsMeta(rootDir)
	if err != nil {
		return nil
	}

	appMetas, infraMetas, monitoringMetas := metas.GetMetasSeparatedByCategories()

	appKustomizations := []any{}
	infraKustomizations := []any{}
	monitoringKustomizations := []any{}

	var wg sync.WaitGroup
	wg.Go(func() {
		sort.Slice(appMetas, func(i, j int) bool {
			return appMetas[i].Name < appMetas[j].Name
		})
		for _, meta := range appMetas {
			appKustomizations = append(appKustomizations, metaToFluxKustomizationManifest(meta))
		}
	})
	wg.Go(func() {
		sort.Slice(infraMetas, func(i, j int) bool {
			return infraMetas[i].Name < infraMetas[j].Name
		})
		for _, meta := range infraMetas {
			infraKustomizations = append(infraKustomizations, metaToFluxKustomizationManifest(meta))
		}
	})
	wg.Go(func() {
		sort.Slice(monitoringMetas, func(i, j int) bool {
			return monitoringMetas[i].Name < monitoringMetas[j].Name
		})
		for _, meta := range monitoringMetas {
			monitoringKustomizations = append(monitoringKustomizations, metaToFluxKustomizationManifest(meta))
		}
	})
	wg.Wait()

	appManifestsConfig := utils.ManifestConfig{
		Filename:  "apps.yaml",
		Manifests: appKustomizations,
	}
	infraManifestsConfig := utils.ManifestConfig{
		Filename:  "infrastructure.yaml",
		Manifests: infraKustomizations,
	}
	monitoringManifestsConfig := utils.ManifestConfig{
		Filename:  "monitoring.yaml",
		Manifests: monitoringKustomizations,
	}

	return utils.MarshalManifests([]utils.ManifestConfig{appManifestsConfig, infraManifestsConfig, monitoringManifestsConfig})
}
