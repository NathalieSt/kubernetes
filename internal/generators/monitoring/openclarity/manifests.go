package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/flux/oci"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
	"path"
)

func createOpenclarityManifests(rootDir string, generatorMeta generator.GeneratorMeta) map[string][]byte {
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

	postgresMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/postgres/main-cluster"))
	if err != nil {
		fmt.Println("An error happened while getting postgres meta for main-cluster")
		return nil
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
						"postgresql": map[string]any{
							"enabled": false,
						},
						"apiserver": map[string]any{
							"database": map[string]any{
								"externalPostgresql": map[string]any{
									"enabled": true,
									"host":    postgresMeta.ClusterUrl,
									"port":    postgresMeta.Port,
									"auth": map[string]any{
										"existingSecret": shared.MainPostgresOpenClarityCredsSecret,
									},
								},
								"postgresql": map[string]any{
									"enabled": false,
								},
							},
						},
						"orchestrator": map[string]any{
							"provider": "kubernetes",
							"serviceAccount": map[string]any{
								"automountServiceAccountToken": true,
							},
						},
						"crDiscoveryServer": map[string]any{
							"env": []core.Env{
								{
									Name:  "CONTAINERD_SOCK_ADDRESS",
									Value: "/run/k0s/containerd.sock",
								},
							},
							"containerRuntimePaths": []map[string]any{
								{
									"name":     "k0s-containerd",
									"path":     "/run/k0s/containerd",
									"readOnly": true,
								},
								{
									"name":     "k0s-containerd-sock",
									"path":     "/run/k0s/containerd.sock",
									"readOnly": true,
								},
							},
						},
					},
				}),
		},
	}

	networkPolicy := utils.ManifestConfig{
		Filename: "network-policy.yaml",
		Manifests: []any{
			networking.NewNetworkPolicy(meta.ObjectMeta{
				Name: fmt.Sprintf("%v-networkpolicy", generatorMeta.Name),
			}, networking.NetworkPolicySpec{
				PolicyTypes: []networking.NetworkPolicyType{networking.Ingress},
				PodSelector: meta.LabelSelector{
					MatchLabels: map[string]string{
						"app.kubernetes.io/name": "openclarity",
					},
				},
				Ingress: []networking.NetworkPolicyIngressRule{
					{
						From: []networking.NetworkPolicyPeer{
							{
								PodSelector: meta.LabelSelector{},
							},
							{
								PodSelector: meta.LabelSelector{
									MachExpressions: []meta.MatchExpression{
										{
											Key:      "app.kubernetes.io/name",
											Operator: meta.In,
											Values: []string{
												"caddy",
											},
										},
									},
								},
								NamespaceSelector: meta.LabelSelector{
									MachExpressions: []meta.MatchExpression{
										{
											Key:      "kubernetes.io/metadata.name",
											Operator: meta.In,
											Values: []string{
												"caddy",
											},
										},
									},
								},
							},
						},
					},
				},
			}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
			release.Filename,
			networkPolicy.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, release, networkPolicy})
}
