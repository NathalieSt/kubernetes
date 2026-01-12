package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
	"path"
)

func createAffineManifests(rootDir string, generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	confPVCName := "config-pvc"
	confPVC := utils.ManifestConfig{
		Filename: "conf-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: confPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			},
			),
		},
	}

	dataPVCName := "data-pvc"
	dataPVC := utils.ManifestConfig{
		Filename: "data-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: dataPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "100Gi",
				}},
				StorageClassName: shared.NFSLocalClass,
			},
			),
		},
	}

	postgresMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/postgres/affine-cluster"))
	if err != nil {
		fmt.Println("An error happened while getting postgres meta ")
		return nil
	}

	entrypointConfigMapName := fmt.Sprintf("%v-configmap", generatorMeta.Name)
	entrypointConfigMap := utils.ManifestConfig{
		Filename: "configmap.yaml",
		Manifests: []any{
			core.NewConfigMap(meta.ObjectMeta{
				Name: entrypointConfigMapName,
				Annotations: map[string]string{
					"app.kubernetes.io/name": generatorMeta.Name,
				},
			}, map[string]string{
				"entrypoint.sh": fmt.Sprintf(`
#!/bin/sh
echo "Setting DATABASE_URL"
export DATABASE_URL="postgresql://${DB_USERNAME}:${DB_PASSWORD}@%v:%v/${DB_DATABASE}"
echo "Done!"
exec "$@"
`, postgresMeta.ClusterUrl, postgresMeta.Port),
			},
			),
		},
	}

	confVolume := "conf-volume"
	dataVolume := "data-volume"
	entrypointVolume := "entrypoint-volume"
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
									Name:    "affine-database-migration",
									Image:   fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Command: []string{"/bin/sh", "-c", "/scripts/entrypoint.sh node ./scripts/self-host-predeploy.js"},
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/root/.affine/config",
											Name:      confVolume,
										},
										{
											MountPath: "/root/.affine/storage",
											Name:      dataVolume,
										},
										{
											MountPath: "/scripts/entrypoint.sh",
											SubPath:   "entrypoint.sh",
											Name:      entrypointVolume,
										},
									},
									Env: []core.Env{
										{
											Name:  "REDIS_SERVER_HOST",
											Value: "redis.redis.svc.cluster.local",
										},
										{
											Name: "DB_USERNAME",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "username",
													Name: shared.AffinePGCredsSecret,
												},
											},
										},
										{
											Name: "DB_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "password",
													Name: shared.AffinePGCredsSecret,
												},
											},
										},
										{
											Name:  "DB_DATABASE",
											Value: "affine-db",
										},
									},
								},
							},
							Containers: []core.Container{
								{
									Name:    generatorMeta.Name,
									Image:   fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Command: []string{"/bin/sh", "-c", "/scripts/entrypoint.sh node ./dist/main.js"},
									Ports: []core.Port{
										{
											ContainerPort: generatorMeta.Port,
											Name:          "affine-ui",
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/root/.affine/config",
											Name:      confVolume,
										},
										{
											MountPath: "/root/.affine/storage",
											Name:      dataVolume,
										},
										{
											MountPath: "/scripts/entrypoint.sh",
											SubPath:   "entrypoint.sh",
											Name:      entrypointVolume,
										},
									},
									Env: []core.Env{
										{
											Name:  "REDIS_SERVER_HOST",
											Value: "redis.redis.svc.cluster.local",
										},
										{
											Name: "DB_USERNAME",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "username",
													Name: shared.AffinePGCredsSecret,
												},
											},
										},
										{
											Name: "DB_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "password",
													Name: shared.AffinePGCredsSecret,
												},
											},
										},
										{
											Name:  "DB_DATABASE",
											Value: "affine-db",
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: confVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: confPVCName,
									},
								},
								{
									Name: dataVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: dataPVCName,
									},
								},
								{
									Name: entrypointVolume,
									ConfigMap: core.ConfigMapVolumeSource{
										Name:        entrypointConfigMapName,
										DefaultMode: 0777,
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

	networkPolicy := utils.ManifestConfig{
		Filename: "network-policy.yaml",
		Manifests: []any{
			networking.NewNetworkPolicy(meta.ObjectMeta{
				Name: fmt.Sprintf("%v-networkpolicy", generatorMeta.Name),
			}, networking.NetworkPolicySpec{
				PolicyTypes: []networking.NetworkPolicyType{networking.Ingress},
				Ingress: []networking.NetworkPolicyIngressRule{
					{
						From: []networking.NetworkPolicyPeer{
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "caddy",
									},
								},
								NamespaceSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"kubernetes.io/metadata.name": "caddy",
									},
								},
							},
						},
					},
				},
			}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			confPVC.Filename,
			dataPVC.Filename,
			service.Filename,
			deployment.Filename,
			entrypointConfigMap.Filename,
			networkPolicy.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, dataPVC, confPVC, kustomization, deployment, service, entrypointConfigMap, networkPolicy})
}
