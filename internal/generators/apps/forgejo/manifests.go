package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/flux/oci"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
	"path"
)

func createForgejoManifests(generatorMeta generator.GeneratorMeta, rootDir string) (map[string][]byte, error) {
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

	postgresMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/shared/infrastructure/postgres/forgejo-cluster"))
	if err != nil {
		fmt.Println("An error happened while getting postgres meta for forgejo")
		return nil, err
	}

	valkeyMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/shared/infrastructure/valkey"))
	if err != nil {
		fmt.Println("An error happened while getting valkey meta for forgejo")
		return nil, err
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
					ValuesFrom: []helm.ReleaseValuesFrom{
						{
							Kind:       helm.Secret,
							Name:       shared.ForgejoPGCredsSecret,
							ValuesKey:  "username",
							TargetPath: "gitea.config.database.USER",
							Optional:   false,
						},
						{
							Kind:       helm.Secret,
							Name:       shared.ForgejoPGCredsSecret,
							ValuesKey:  "password",
							TargetPath: "gitea.config.database.PASSWD",
							Optional:   false,
						},
					},
					Values: map[string]any{
						"gitea": map[string]any{
							"admin": map[string]any{"username": "Nathi"},
							"config": map[string]any{
								"database": map[string]any{
									"DB_TYPE": "postgres",
									"HOST":    fmt.Sprintf("%v:%v", postgresMeta.ClusterUrl, postgresMeta.Port),
									"NAME":    "forgejo",
								},
								"server": map[string]any{
									"ROOT_URL": fmt.Sprintf("https://%v.%v", generatorMeta.Caddy.DNSName, shared.NetbirdDomainBase),
								},
							},
							"queue": map[string]any{
								"TYPE":     "redis",
								"CONN_STR": fmt.Sprintf("valkey://%v:%v/0?", valkeyMeta.ClusterUrl, valkeyMeta.Port),
							},
							"cache": map[string]any{
								"ADAPTER": "redis",
								"HOST":    fmt.Sprintf("valkey://%v:%v/1", valkeyMeta.ClusterUrl, valkeyMeta.Port),
							},
							"session": map[string]any{
								"PROVIDER":        "redis",
								"PROVIDER_CONFIG": fmt.Sprintf("valkey://%v:%v/2", valkeyMeta.ClusterUrl, valkeyMeta.Port),
							},
						},
						"persistence": map[string]any{
							"enabled":      true,
							"storageClass": shared.NFSRemoteClass,
						},
						"service": map[string]any{
							"http": map[string]any{
								"annotations": map[string]any{
									"netbird.io/expose": "true",
									"netbird.io/groups": "cluster-services",
								},
							},
						},
					},
				}),
		},
	}

	scaledObject := utils.ManifestConfig{
		Filename:  "scaled-object.yaml",
		Manifests: utils.GenerateCronScaler(fmt.Sprintf("%v-scaledobject", generatorMeta.Name), generatorMeta.Name, keda.Deployment, generatorMeta.KedaScaling),
	}

	networkPolicy := utils.ManifestConfig{
		Filename: "network-policy.yaml",
		Manifests: []any{
			networking.NewNetworkPolicy(meta.ObjectMeta{
				Name: fmt.Sprintf("%v-networkpolicy", generatorMeta.Name),
			}, networking.NetworkPolicySpec{
				PolicyTypes: []networking.NetworkPolicyType{networking.Ingress},
				Ingress: []networking.NetworkPolicyIngressRule{
					{
						From: []networking.NetworkPolicyPeer{
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "caddy",
									},
								},
								NamespaceSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"kubernetes.io/metadata.name": "caddy",
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
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				namespace.Filename,
				repo.Filename,
				release.Filename,
				scaledObject.Filename,
				networkPolicy.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, release, scaledObject, networkPolicy}), nil
}
