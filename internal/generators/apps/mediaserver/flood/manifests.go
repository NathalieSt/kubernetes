package main

import (
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createTransmissionManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	floodConfigPVCName := "flood-config-pvc"
	floodConfigPVC := utils.ManifestConfig{
		Filename: "flood-config-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: floodConfigPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSLocalClass,
			}),
		},
	}

	mediaVolume := "media-volume"
	floodConfigVolume := "flood-config-volume"
	deployment := utils.ManifestConfig{
		Filename: "deployment.yaml",
		Manifests: []any{
			apps.NewDeployment(
				meta.ObjectMeta{
					Name: "flood",
					Labels: map[string]string{
						"app.kubernetes.io/name":    "flood",
						"app.kubernetes.io/version": "1.0",
					},
				},
				apps.DeploymentSpec{
					Replicas: 1,
					Selector: meta.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":    "flood",
							"app.kubernetes.io/version": "1.0",
						},
					},
					Template: core.PodTemplateSpec{
						Metadata: meta.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":    "flood",
								"app.kubernetes.io/version": "1.0",
							},
						},
						Spec: core.PodSpec{
							Containers: []core.Container{
								{
									Name:  "flood-ui",
									Image: "jesec/flood:4.11",
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/data",
											Name:      mediaVolume,
										},
										{
											MountPath: "/config",
											Name:      floodConfigVolume,
										},
									},
									Env: []core.Env{
										{
											Name:  "HOME",
											Value: "/config",
										},
									},
									Ports: []core.Port{
										{
											Name:          "flood-web-ui",
											ContainerPort: 3000,
										},
									},
								},
							},
							Volumes: []core.Volume{

								{
									Name: floodConfigVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: floodConfigPVCName,
									},
								},
								{
									Name: mediaVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										//FIXME: get from jellyfin generator
										ClaimName: "jellyfin-pvc",
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
					Name: "flood",
					Labels: map[string]string{
						"app.kubernetes.io/name":    "flood",
						"app.kubernetes.io/version": "1.0",
					},
				}, core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name":    "flood",
						"app.kubernetes.io/version": "1.0",
					},
					Ports: []core.ServicePort{
						{
							Name:       "http-flood-webui",
							Port:       3000,
							TargetPort: 3000,
						},
					},
				},
			),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			deployment.Filename,
			service.Filename,
			floodConfigPVC.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, deployment, service, floodConfigPVC})
}
