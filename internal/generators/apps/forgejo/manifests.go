package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/generators/infrastructure"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/flux/oci"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/kustomize"
)

func createForgejoManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename: "namespace.yaml",
		Generate: func() any {
			return core.NewNamespace(meta.ObjectMeta{
				Name: generatorMeta.Namespace,
				Labels: map[string]string{
					"istio-injection": "enabled",
				},
			})
		},
	}

	repoName := fmt.Sprintf("%v-repo", generatorMeta.Name)
	repo := utils.ManifestConfig{
		Filename: "repo.yaml",
		Generate: func() any {
			return oci.NewRepo(meta.ObjectMeta{
				Name: repoName,
			}, oci.RepoSpec{
				Url:      generatorMeta.Helm.Url,
				Interval: "24h",
				Ref: oci.RepoRef{
					Tag: generatorMeta.Helm.Version,
				},
			})
		},
	}

	release := utils.ManifestConfig{
		Filename: "release.yaml",
		Generate: func() any {
			return helm.NewRelease(meta.ObjectMeta{
				Name: generatorMeta.Name,
			}, helm.ReleaseSpec{
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
						Name:       generators.PostgresCredsSecret,
						ValuesKey:  "username",
						TargetPath: "gitea.config.database.USER",
						Optional:   false,
					},
					{
						Kind:       helm.Secret,
						Name:       generators.PostgresCredsSecret,
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
								"HOST":    fmt.Sprintf("%v:%v", infrastructure.Postgres.ClusterUrl, infrastructure.Postgres.Port),
								"NAME":    "forgejo",
							},
						},
						"queue": map[string]any{
							"TYPE":     "redis",
							"CONN_STR": "valkey://valkey.valkey.svc.cluster.local:6379/0?",
						},
						"cache": map[string]any{
							"ADAPTER": "redis",
							"HOST":    "valkey://valkey.valkey.svc.cluster.local:6379/1",
						},
						"session": map[string]any{
							"PROVIDER":        "redis",
							"PROVIDER_CONFIG": "valkey://valkey.valkey.svc.cluster.local:6379/2",
						},
						"podAnnotations": map[string]any{
							"proxy.istio.io/config": "{'holdApplicationUntilProxyStarts': true}",
						},
					},
					"persistence": map[string]any{
						"enabled":      true,
						"storageClass": generators.NFSRemoteClass,
					},
				},
			})
		},
	}

	scaledObject := utils.ManifestConfig{
		Filename: "scaled-object.yaml",
		Generate: func() any {
			return keda.NewScaledObject(
				meta.ObjectMeta{
					Name: fmt.Sprintf("%v-scaledobject", generatorMeta.Name),
				}, keda.ScaledObjectSpec{
					ScaleTargetRef: keda.ScaleTargetRef{
						Name: generatorMeta.Name,
					},
					MinReplicaCount: 0,
					CooldownPeriod:  300,
					Triggers: []keda.ScaledObjectTrigger{
						{
							ScalerType: keda.Cron,
							Metadata:   generatorMeta.KedaScaling,
						},
					},
				},
			)
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Generate: func() any {
			return kustomize.NewKustomization(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
				},
				[]string{
					namespace.Filename,
					repo.Filename,
					release.Filename,
					scaledObject.Filename,
				},
			)
		},
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, release, scaledObject})
}
