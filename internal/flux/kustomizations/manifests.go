package main

import (
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/kustomization"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
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

	print(metas)
	appMetas, infraMetas, monitoringMetas := metas.GetMetasSeparatedByCategories()

	appManifests := []any{}
	infraManifests := []any{}
	monitoringManifests := []any{}

	var wg sync.WaitGroup
	wg.Go(func() {
		for _, meta := range appMetas {
			appManifests = append(appManifests, metaToFluxKustomizationManifest(meta))
		}
	})
	wg.Go(func() {
		for _, meta := range infraMetas {
			infraManifests = append(appManifests, metaToFluxKustomizationManifest(meta))
		}
	})
	wg.Go(func() {
		for _, meta := range monitoringMetas {
			monitoringManifests = append(appManifests, metaToFluxKustomizationManifest(meta))
		}
	})
	wg.Wait()

	appManifestsConfig := utils.ManifestConfig{
		Filename:  "apps.yaml",
		Manifests: appManifests,
	}
	infraManifestsConfig := utils.ManifestConfig{
		Filename:  "infrastructure.yaml",
		Manifests: infraManifests,
	}
	monitoringManifestsConfig := utils.ManifestConfig{
		Filename:  "monitoring.yaml",
		Manifests: monitoringManifests,
	}

	return utils.MarshalManifests([]utils.ManifestConfig{appManifestsConfig, infraManifestsConfig, monitoringManifestsConfig})
}
