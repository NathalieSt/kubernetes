package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/infrastructure/vaultsecretsoperator"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createVaultSecretsOperatorManifests(generatorMeta generator.GeneratorMeta, rootDir string) (map[string][]byte, error) {
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
						"controller": map[string]any{
							"annotations": map[string]any{
								"traffic.sidecar.istio.io/excludeOutboundPorts": "8200",
							},
						},
					},
				}),
		},
	}

	vaultMeta, err := utils.GetServiceMeta(rootDir, "internal/generators/infrastructure/vault")
	if err != nil {
		fmt.Println("An error happened while getting vault meta ")
		return nil, err
	}

	serviceAccount, role, rolebinding := utils.GenerateRBAC(generatorMeta.Name)

	vaultConfigs := utils.ManifestConfig{
		Filename: "vault-configs.yaml",
		Manifests: []any{
			vaultsecretsoperator.NewAuthGlobal(meta.ObjectMeta{
				Name: "default",
			}, vaultsecretsoperator.AuthGlobalSpec{
				AllowedNamespaces: []string{"reflector", "gluetun-proxy"},
				DefaultAuthMethod: "kubernetes",
				Kubernetes: vaultsecretsoperator.Kubernetes{
					Audiences:              []string{"vault"},
					Mount:                  "kubernetes",
					Role:                   "global-vault-auth",
					ServiceAccount:         serviceAccount.Metadata.Name,
					TokenExpirationSeconds: 600,
				},
			}),
			serviceAccount,
			role,
			rolebinding,
			vaultsecretsoperator.NewConnection(meta.ObjectMeta{
				Name: "default",
			}, vaultsecretsoperator.ConnectionSpec{
				Address: fmt.Sprintf("http://%v:8200", vaultMeta.ClusterUrl),
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
				vaultConfigs.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, vaultConfigs}), nil
}
