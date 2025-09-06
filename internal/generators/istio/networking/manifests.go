package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createJellyfinManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename: "namespace.yaml",
		Manifests: []any{
			core.NewNamespace(meta.ObjectMeta{
				Name: generatorMeta.Namespace,
				Labels: map[string]string{
					"istio-injection": "enabled",
				},
			}),
		},
	}

	pvcName := fmt.Sprintf("%v-pvc", generatorMeta.Name)
	pvc := utils.ManifestConfig{
		Filename: "pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "100Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			}),
		},
	}

	repoName := fmt.Sprintf("%v-repo", generatorMeta.Name)
	repo := utils.ManifestConfig{
		Filename: "repo.yaml",
		Manifests: []any{
			helm.NewRepo(meta.ObjectMeta{
				Name: repoName,
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
			}, helm.ReleaseSpec{
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
					"persistence": map[string]any{
						"media": map[string]any{
							"existingClaim": pvcName,
						},
					},
				},
			}),
		},
	}

	deploymentName := "jellyfin"
	scaledObject := utils.ManifestConfig{
		Filename: "scaled-object.yaml",
		Manifests: []any{
			keda.NewScaledObject(
				meta.ObjectMeta{
					Name: fmt.Sprintf("%v-scaledobject", generatorMeta.Name),
				}, keda.ScaledObjectSpec{
					ScaleTargetRef: keda.ScaleTargetRef{
						Name: deploymentName,
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
			),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
			chart.Filename,
			release.Filename,
			pvc.Filename,
			scaledObject.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, pvc, scaledObject})
}
