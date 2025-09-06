package main

import (
	"fmt"
	"kubernetes/internal/generators/istio"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/kustomize"
)

func createBaseManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	chartName := fmt.Sprintf("%v-chart", generatorMeta.Name)
	chart := utils.ManifestConfig{
		Filename: "chart.yaml",
		Manifests: []any{
			helm.NewChart(
				meta.ObjectMeta{
					Name: chartName,
				},
				helm.ChartSpec{
					Interval:          "24h",
					Chart:             generatorMeta.Helm.Chart,
					ReconcileStrategy: helm.ChartVersion,
					SourceRef: helm.ChartSourceRef{
						Kind: helm.HelmRepository,
						Name: istio.RepoName,
					},
					Version: generatorMeta.Helm.Version,
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
						Kind: helm.HelmChart,
						Name: chartName,
					},
					Install: helm.ReleaseInstall{
						Remediation: helm.ReleaseInstallRemediation{
							Retries: 3,
						},
					},
				}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: []any{
			kustomize.NewKustomization(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
				},
				[]string{
					release.Filename,
				},
			),
		},
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, chart, release})
}
