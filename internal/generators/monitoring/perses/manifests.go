package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
)

func createPersesManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"persistence": map[string]any{
				"enabled":      true,
				"storageClass": generators.NFSRemoteClass,
			},
			"provisioningPersistence": map[string]any{
				"enabled":      true,
				"storageClass": generators.NFSRemoteClass,
			},
			"service": map[string]any{
				"annotations": map[string]any{
					"netbird.io/expose": "true",
					"netbird.io/groups": "cluster-services",
				},
			},
		},
		nil,
	)

	scaledObject := utils.ManifestConfig{
		Filename:  "scaled-object.yaml",
		Manifests: utils.GenerateCronScaler(fmt.Sprintf("%v-scaledobject", generatorMeta.Name), generatorMeta.Name, keda.StatefulSet, generatorMeta.KedaScaling),
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
			chart.Filename,
			release.Filename,
			scaledObject.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, scaledObject})
}
