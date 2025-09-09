package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createVaultManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	repoName := fmt.Sprintf("%v-repo", generatorMeta.Name)
	repo := utils.ManifestConfig{
		Filename: "repo.yaml",
		Manifests: []any{
			helm.NewRepo(meta.ObjectMeta{
				Name: repoName,
			},
				helm.RepoSpec{
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
			helm.NewChart(meta.ObjectMeta{
				Name: chartName,
			}, helm.ChartSpec{
				Interval:          "24h",
				Chart:             generatorMeta.Helm.Chart,
				ReconcileStrategy: helm.ChartVersion,
				SourceRef: helm.ChartSourceRef{
					Kind: helm.HelmRepository,
					Name: repoName,
				},
				Version: generatorMeta.Helm.Version,
			}),
		},
	}

	release := utils.ManifestConfig{
		Filename: "release.yaml",
		Manifests: []any{
			helm.NewRelease(meta.ObjectMeta{
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
						"ui": map[string]any{
							"enabled": true,
						},
						"server": map[string]any{
							"dataStorage": map[string]any{
								"storageClass": generators.NFSLocalClass,
							},
						},
					},
				}),
		},
	}

	peerAuth := utils.ManifestConfig{
		Filename: "peer-auth.yaml",
		Manifests: []any{
			istio.NewPeerAuthenthication(meta.ObjectMeta{
				Name: "vault-permissive-mtls",
			}, istio.PeerAuthenthicationSpec{
				MTLS: istio.PeerAuthenthicationmTLS{
					Mode: istio.PERMISSIVE,
				},
			}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				namespace.Filename,
				repo.Filename,
				chart.Filename,
				release.Filename,
				peerAuth.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, peerAuth})
}
