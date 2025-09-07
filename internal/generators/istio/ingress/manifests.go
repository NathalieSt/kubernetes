package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/generators/istio"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/kustomize"
)

func createIngressManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
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
				},
			),
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
					Values: map[string]any{
						"autoscaling": map[string]any{
							"maxReplicas": 2,
						},
						"service": map[string]any{
							"loadBalancerIP": generators.IstioGatewayIP,
						},
					},
				},
			),
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
