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

	datapvcName := "data-pvc"
	datapvc := utils.ManifestConfig{
		Filename: "data-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: datapvcName,
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

	bookspvcName := "books-pvc"
	bookspvc := utils.ManifestConfig{
		Filename: "books-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: bookspvcName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "100Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			},
			),
		},
	}

	bookdroppvcName := "bookdroppvc-pvc"
	bookdroppvc := utils.ManifestConfig{
		Filename: "bookdrop-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: bookdroppvcName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "10Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			},
			),
		},
	}

	dataVolume := "data-volume"
	booksVolume := "books-volume"
	bookdropVolume := "bookdrop-volume"
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
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/app/data",
											Name:      dataVolume,
										},
										{
											MountPath: "/books",
											Name:      booksVolume,
										},
										{
											MountPath: "/bookdrop",
											Name:      bookdropVolume,
										},
									},
									Env: []core.Env{
										{
											Name:  "USER_ID",
											Value: "0",
										},
										{
											Name:  "GROUP_ID",
											Value: "0",
										},
										{
											Name:  "TZ",
											Value: "Europe/Vienna",
										},
										{
											Name:  "DATABASE_URL",
											Value: "jdbc:mariadb://mariadb.mariadb.svc.cluster.local:3306/booklore",
										},
										{
											Name: "DATABASE_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "password",
													Name: generators.MariaDBCredsSecret,
												},
											},
										},
										{
											Name: "DATABASE_USERNAME",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "username",
													Name: generators.MariaDBCredsSecret,
												},
											},
										},
										{
											Name:  "BOOKLORE_PORT",
											Value: "6060",
										},
										{
											Name:  "SWAGGER_ENABLED",
											Value: "false",
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: dataVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: datapvcName,
									},
								},
								{
									Name: booksVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: bookspvcName,
									},
								},
								{
									Name: bookdropVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: bookdroppvcName,
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

	scaledObject := utils.ManifestConfig{
		Filename:  "scaled-object.yaml",
		Manifests: utils.GenerateCronScaler(fmt.Sprintf("%v-scaledobject", generatorMeta.Name), generatorMeta.Name, generatorMeta.KedaScaling),
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			bookdroppvc.Filename,
			bookspvc.Filename,
			datapvc.Filename,
			deployment.Filename,
			service.Filename,
			scaledObject.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, bookdroppvc, bookspvc, datapvc, kustomization, deployment, service, scaledObject})
}
