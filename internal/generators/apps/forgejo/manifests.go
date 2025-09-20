package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/flux/oci"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"path"
)

func createForgejoManifests(generatorMeta generator.GeneratorMeta, rootDir string) (map[string][]byte, error) {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
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

	postgresMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/postgres/forgejo-cluster"))
	if err != nil {
		fmt.Println("An error happened while getting postgres meta for forgejo")
		return nil, err
	}

	valkeyMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/valkey"))
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
							Name:       generators.ForgejoPGCredsSecret,
							ValuesKey:  "username",
							TargetPath: "gitea.config.database.USER",
							Optional:   false,
						},
						{
							Kind:       helm.Secret,
							Name:       generators.ForgejoPGCredsSecret,
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
									"ROOT_URL": fmt.Sprintf("https://%v.netbird.selfhosted", generatorMeta.Caddy.DNSName),
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
							"storageClass": generators.NFSRemoteClass,
						},
					},
				}),
		},
	}

	scaledObject := utils.ManifestConfig{
		Filename:  "scaled-object.yaml",
		Manifests: utils.GenerateCronScaler(fmt.Sprintf("%v-scaledobject", generatorMeta.Name), generatorMeta.Name, generatorMeta.KedaScaling),
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
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, release, scaledObject}), nil
}
