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

func createInvidiousManifests(generatorMeta generator.GeneratorMeta, rootDir string, relativeDir string) (map[string][]byte, error) {
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

	configMapName := "invidious-cm"
	configMap, err := getInvidiousConfigMap(rootDir, relativeDir, configMapName)
	if err != nil {
		fmt.Println("An error occurred while getting the configMap for discord-bridge")
		return nil, err
	}

	configMapManifest := utils.ManifestConfig{
		Filename:  "configmap.yaml",
		Manifests: []any{*configMap},
	}

	configPVCName := "config-pvc"
	configPVC := utils.ManifestConfig{
		Filename: "config-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: configPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: generators.NFSRemoteClass,
			},
			),
		},
	}

	configPVCVolumeName := "config-pvc-volume"
	configMapVolumeName := "configmap-volume"
	cachePVCVolume := "cache-pvc-volume"
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
									Name:  "config-init",
									Image: "alpine:latest",
									Command: []string{
										"/bin/sh",
										"-c",
										`apk update && apk add gettext;
envsubst < /template/config.yaml > /data/config.yml;
										`,
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      configMapVolumeName,
											MountPath: "/template",
										},
										{
											Name:      configPVCVolumeName,
											MountPath: "/data",
										},
									},
									Env: []core.Env{
										{
											Name:  "PG_USER",
											Value: "Test",
										},
										{
											Name:  "PG_PASSWORD",
											Value: "Test stuff",
										},

										{
											Name:  "HMAC_KEY",
											Value: "1234567890123456",
										},
										{
											Name:  "COMPANION_KEY",
											Value: "1234567890123456",
										},
									},
								},
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
											MountPath: "/config/",
											Name:      configPVCVolumeName,
										},
									},
								},
								{
									Name:  "invidious-companion",
									Image: "quay.io/invidious/invidious-companion:latest",
									Env:   []core.Env{},
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/var/tmp/youtubei.js",
											Name:      configMapVolumeName,
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: configMapVolumeName,
									ConfigMap: core.ConfigMapVolumeSource{
										Name: configMapName,
									},
								},
								{
									Name: configPVCVolumeName,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: configPVCName,
									},
								},
								{
									Name: cachePVCVolume,
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
			configMapManifest.Filename,
			configPVC.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, cachePVC, deployment, service, configMapManifest, configPVC}), nil
}
