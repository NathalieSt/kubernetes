package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/flux/oci"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createVictoriaMetricsManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repoName := fmt.Sprintf("%v-repo", generatorMeta.Name)
	repo := utils.ManifestConfig{
		Filename: "repo.yaml",
		Manifests: []any{
			oci.NewRepo(
				meta.ObjectMeta{
					Name: repoName,
				},
				oci.RepoSpec{
					Url:      generatorMeta.Helm.Url,
					Interval: "24h",
					Ref: oci.RepoRef{
						Tag: generatorMeta.Helm.Version,
					},
				}),
		},
	}

	release := utils.ManifestConfig{
		Filename: "release.yaml",
		Manifests: []any{
			helm.NewRelease(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
				},
				helm.ReleaseSpec{
					ReleaseName: generatorMeta.Name,
					Interval:    "24h",
					Timeout:     "10m",
					ChartRef: helm.ReleaseChartRef{
						Kind: helm.OCIRepository,
						Name: repoName,
					},
					Install: helm.ReleaseInstall{
						Remediation: helm.ReleaseInstallRemediation{
							Retries: 3,
						},
					},
					Values: map[string]any{
						"nameOverride":      "vmks",
						"grafana":           map[string]any{"enabled": false},
						"defaultDashboards": map[string]any{"enabled": false},
						"vmsingle": map[string]any{
							"spec": map[string]any{
								"storage": map[string]any{
									"storageClassName": generators.NFSRemoteClass,
								},
							},
						},
					},
				},
			),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
			release.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, release})
}
