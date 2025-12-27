package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createLidifyManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	musicPVCName := fmt.Sprintf("%v-music-pvc", generatorMeta.Name)
	musicPVC := utils.ManifestConfig{
		Filename: "music-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: musicPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "100Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			}),
		},
	}

	configPVCName := fmt.Sprintf("%v-config-pvc", generatorMeta.Name)
	configPVC := utils.ManifestConfig{
		Filename: "config-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: configPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteOnce"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "30Gi",
				}},
				StorageClassName: "local-path",
			}),
		},
	}

	configVolumeName := "config-volume"
	musicVolumeName := "music-volume"
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
							InitContainers: []core.Container{
								{
									Name:    "init-redis-datadir",
									Command: []string{"sh", "-c"},
									Args:    []string{"mkdir -p /data/redis"},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      configVolumeName,
											MountPath: "/data",
										},
									},
								},
							},
							// use while no arm image is available
							NodeSelector: map[string]string{
								"kubernetes.io/hostname": "debian",
							},
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
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/music",
											Name:      musicVolumeName,
										},
										{
											MountPath: "/data",
											Name:      configVolumeName,
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: configVolumeName,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: configPVCName,
									},
								}, {
									Name: musicVolumeName,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: musicPVCName,
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
					Annotations: map[string]string{
						"netbird.io/expose": "true",
						"netbird.io/groups": "cluster-services",
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

	scaledObject := utils.ManifestConfig{
		Filename:  "scaled-object.yaml",
		Manifests: utils.GenerateCronScaler(fmt.Sprintf("%v-scaledobject", generatorMeta.Name), generatorMeta.Name, keda.Deployment, generatorMeta.KedaScaling),
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			deployment.Filename,
			service.Filename,
			configPVC.Filename,
			musicPVC.Filename,
			scaledObject.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment, service, configPVC, musicPVC, scaledObject})
}
