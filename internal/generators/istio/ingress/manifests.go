package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/kustomize"
)

func createIngressManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"autoscaling": map[string]any{
				"maxReplicas": 2,
			},
			"service": map[string]any{
				"loadBalancerIP": generators.IstioGatewayIP,
			},
		},
		nil,
	)

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: []any{
			kustomize.NewKustomization(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
				},
				[]string{
					repo.Filename,
					chart.Filename,
					release.Filename,
					namespace.Filename,
				},
			),
		},
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, repo, kustomization, chart, release})
}
