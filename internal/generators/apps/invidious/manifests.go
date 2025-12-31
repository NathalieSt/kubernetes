package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createBookloreManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	cachePVCName := "cache-pvc"
	cachePVC := utils.ManifestConfig{
		Filename: "cache-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: cachePVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "10Gi",
				}},
				StorageClassName: generators.NFSRemoteClass,
			},
			),
		},
	}

	cacheVolume := "cache-volume"
	deployment := utils.ManifestConfig{
		Filename: "deployment.yaml",
		Manifests: []any{
			apps.NewDeployment(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				},
				apps.DeploymentSpec{
					Replicas: 1,
					Selector: meta.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name": generatorMeta.Name,
						},
					},
					Template: core.PodTemplateSpec{
						Metadata: meta.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":    generatorMeta.Name,
								"app.kubernetes.io/version": generatorMeta.Docker.Version,
							},
						},
						Spec: core.PodSpec{
							Containers: []core.Container{
								{
									Name:  generatorMeta.Name,
									Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Ports: []core.Port{
										{
											ContainerPort: generatorMeta.Port,
											Name:          generatorMeta.Name,
										},
									},
									Env: []core.Env{
										core.Env{
											Name:      "INVIDIOUS_CONFIG",
											ValueFrom: core.ValueFrom{},
										},
									},
								},
								{
									Name:  "invidious-companion",
									Image: "quay.io/invidious/invidious-companion:latest",
									Env: []core.Env{
										{
											Name: "NB_SETUP_KEY",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: generators.NetbirdSecretName,
													Key:  "setup-key",
												},
											},
										},
										{
											Name:  "NB_HOSTNAME",
											Value: "adguard-home",
										},
										{
											Name:  "NB_MANAGEMENT_URL",
											Value: "https://netbird.nathalie-stiefsohn.eu",
										},
									},
									VolumeMounts: []core.VolumeMount{
										core.VolumeMount{
											MountPath: "/var/tmp/youtubei.js",
											Name:      cacheVolume,
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: cacheVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: cachePVCName,
									},
								},
							},
						},
					},
				},
			),
		},
	}

	service := utils.ManifestConfig{
		Filename: "service.yaml",
		Manifests: []any{
			core.NewService(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				}, core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
					Ports: []core.ServicePort{
						{
							Name:       fmt.Sprintf("http-%v", generatorMeta.Name),
							Port:       generatorMeta.Port,
							TargetPort: generatorMeta.Port,
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
			service.Filename,
			cachePVC.Filename,
			deployment.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, cachePVC, deployment, service})
}
