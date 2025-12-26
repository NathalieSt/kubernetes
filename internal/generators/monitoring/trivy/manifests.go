package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
)

func createKialiManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"alternateReportStorage": map[string]any{
				"storageClassName": generators.DebianStorageClass,
			},
			"trivy": map[string]any{
				"storageClassName": generators.DebianStorageClass,
				"debug":            true,
			},
			"operator": map[string]any{
				"scanJobsConcurrentLimit": 2,
				"logDevMode":              true,
				"scanJobPodTemplateResources": map[string]any{
					"limits": map[string]any{
						"cpu":    "1",
						"memory": "1Gi",
					},
					"requests": map[string]any{
						"cpu":    "500m",
						"memory": "512Mi",
					},
				},
				"scanJobPodTemplateContainerSecurityContext": map[string]any{
					"runAsUser":                0,
					"runAsGroup":               0,
					"readOnlyRootFilesystem":   false,
					"allowPrivilegeEscalation": false,
					"capabilities": map[string]any{
						"drop": []string{"ALL"},
					},
					"privileged": false,
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
