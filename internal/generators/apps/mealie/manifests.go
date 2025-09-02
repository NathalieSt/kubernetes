package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/generators/infrastructure"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/kustomize"
)

func CreateMealieManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
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

	pvcName := fmt.Sprintf("%v-pvc", generatorMeta.Name)
	pvc := utils.ManifestConfig{
		Filename: "pvc.yaml",
		Generate: func() any {
			return core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "100Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			})
		},
	}

	volumeName := "mealie-pvc"
	deployment := utils.ManifestConfig{
		Filename: "deployment.yaml",
		Generate: func() any {
			return apps.NewDeployment(
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
									Name:  fmt.Sprintf("%v-deployment", generatorMeta.Name),
									Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Ports: []core.Port{
										{
											ContainerPort: generatorMeta.Port,
											Name:          generatorMeta.Name,
										},
									},
									Env: []core.Env{
										{
											Name:  "POSTGRES_SERVER",
											Value: infrastructure.Postgres.ClusterUrl,
										},
										{
											Name:  "POSTGRES_PORT",
											Value: fmt.Sprintf("%v", infrastructure.Postgres.Port),
										},
										{
											Name: "POSTGRES_USERNAME",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "username",
													Name: generators.PostgresCredsSecret,
												},
											},
										},
										{
											Name: "POSTGRES_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "password",
													Name: generators.PostgresCredsSecret,
												},
											},
										},
										{
											Name:  "POSTGRES_DB",
											Value: "mealie",
										},
										{
											Name:  "BASE_URL",
											Value: "https://mealie.cluster.netbird.selfhosted",
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/app/data/",
											Name:      volumeName,
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: volumeName,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: pvcName,
									},
								},
							},
						},
					},
				},
			)
		},
	}

	service := utils.ManifestConfig{
		Filename: "service.yaml",
		Generate: func() any {
			return core.NewService(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
				}, core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
					Ports: []core.ServicePort{
						{
							Name:       fmt.Sprintf("http-%v", generatorMeta.Name),
							Port:       9000,
							TargetPort: 9000,
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
					deployment.Filename,
					pvc.Filename,
					service.Filename,
				},
			)
		},
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment, pvc, service})
}
