package main

import (
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createKedaManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			// Pod annotations for KEDA operator
			"podAnnotations": map[string]any{
				"keda": map[string]any{
					"traffic.sidecar.istio.io/excludeInboundPorts":  "9666",
					"traffic.sidecar.istio.io/excludeOutboundPorts": "9443,6443",
				},
				// Pod annotations for KEDA Metrics Adapter
				"metricsAdapter": map[string]any{
					"traffic.sidecar.istio.io/excludeInboundPorts":  "6443",
					"traffic.sidecar.istio.io/excludeOutboundPorts": "9666,9443",
				},
				// Pod annotations for KEDA Admission webhooks
				"webhooks": map[string]any{
					"traffic.sidecar.istio.io/excludeInboundPorts":  "9443",
					"traffic.sidecar.istio.io/excludeOutboundPorts": "9666,6443",
				},
			},
		},
		nil,
	)

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
