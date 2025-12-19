package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"path"
)

func createMealieManifests(generatorMeta generator.GeneratorMeta, rootDir string) (map[string][]byte, error) {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
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
			},
			),
		},
	}

	postgresMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/postgres"))
	if err != nil {
		fmt.Println("An error happened while getting postgres meta ")
		return nil, err
	}

	volumeName := "mealie-pvc"
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
										{
											Name:  "POSTGRES_SERVER",
											Value: postgresMeta.ClusterUrl,
										},
										{
											Name:  "POSTGRES_PORT",
											Value: fmt.Sprintf("%v", postgresMeta.Port),
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
											Value: fmt.Sprintf("https://%v.netbird.selfhosted", generatorMeta.Caddy.DNSName),
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
			deployment.Filename,
			pvc.Filename,
			service.Filename,
			scaledObject.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment, pvc, service, scaledObject}), nil
}
