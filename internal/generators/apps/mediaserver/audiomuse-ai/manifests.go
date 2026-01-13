package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"path"
)

func createAudiomuseAIManifests(rootDir string, generatorMeta generator.GeneratorMeta) map[string][]byte {

	secretName := "jellyfin-creds"
	jellyfinCredsSecret := utils.StaticSecretConfig{
		Name:       secretName,
		SecretName: secretName,
		Path:       "jellyfin",
	}

	jellyfinCredsVaultSecret := utils.ManifestConfig{
		Filename: "vault-secret.yaml",
		Manifests: utils.GenerateVaultAccessManifests(
			generatorMeta.Name,
			//FIXME: get this from VSO generator meta
			"vault-secrets-operator",
			[]utils.StaticSecretConfig{jellyfinCredsSecret},
		),
	}

	postgresMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/postgres/main-cluster"))
	if err != nil {
		fmt.Println("An error happened while getting postgres meta ")
		return nil
	}

	redisMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/redis"))
	if err != nil {
		fmt.Println("An error happened while getting redis meta ")
		return nil
	}

	jellyfinMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/apps/mediaserver/jellyfin"))
	if err != nil {
		fmt.Println("An error happened while getting jellyfin meta ")
		return nil
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
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}

	configMapName := "audiomuse-configmap"
	configMap := utils.ManifestConfig{
		Filename: "configmap.yaml",
		Manifests: []any{
			core.NewConfigMap(meta.ObjectMeta{
				Name: configMapName,
			}, map[string]string{
				"MEDIASERVER_TYPE": "jellyfin",
				"JELLYFIN_URL":     fmt.Sprintf("http://%v:%v", jellyfinMeta.ClusterUrl, jellyfinMeta.Port),
				"POSTGRES_HOST":    postgresMeta.ClusterUrl,
				"POSTGRES_PORT":    fmt.Sprint(postgresMeta.Port),
				"REDIS_URL":        fmt.Sprintf("redis://%v:%v/0", redisMeta.ClusterUrl, redisMeta.Port),
				"CLAP_ENABLED":     "true",
				"TEMP_DIR":         "/app/temp_audio",
				"POSTGRES_DB":      "audiomuseai",
			}),
		},
	}

	tempVolumeName := "temp-volume"

	flaskDeployment := utils.ManifestConfig{
		Filename: "flask-deployment.yaml",
		Manifests: []any{
			apps.NewDeployment(
				meta.ObjectMeta{
					Name: "audiomuse-ai-flask",
					Labels: map[string]string{
						"app.kubernetes.io/name":    fmt.Sprintf("%v-flask", generatorMeta.Name),
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				},
				apps.DeploymentSpec{
					Replicas: 1,
					Selector: meta.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":    fmt.Sprintf("%v-flask", generatorMeta.Name),
							"app.kubernetes.io/version": generatorMeta.Docker.Version,
						},
					},
					Template: core.PodTemplateSpec{
						Metadata: meta.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":    fmt.Sprintf("%v-flask", generatorMeta.Name),
								"app.kubernetes.io/version": generatorMeta.Docker.Version,
							},
						},
						Spec: core.PodSpec{
							Containers: []core.Container{
								{
									Name:  "audiomuse-ai-flask",
									Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Env: []core.Env{
										{
											Name:  "SERVICE_TYPE",
											Value: "flask",
										},
										{
											Name:  "CLAP_ENABLED",
											Value: "false",
										},
										{
											Name: "JELLYFIN_USER_ID",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: secretName,
													Key:  "userid",
												},
											},
										},
										{
											Name: "JELLYFIN_TOKEN",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: secretName,
													Key:  "token",
												},
											},
										},
										{
											Name: "POSTGRES_USER",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: shared.PostgresCredsSecret,
													Key:  "username",
												},
											},
										},
										{
											Name: "POSTGRES_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: shared.PostgresCredsSecret,
													Key:  "password",
												},
											},
										},
									},
									EnvFrom: []core.EnvReference{
										{
											ConfigMapRef: core.ConfigMapRef{
												Name: configMapName,
											},
										},
									},
									Ports: []core.Port{
										{
											Name:          "flask",
											ContainerPort: generatorMeta.Port,
										},
									},
									VolumeMounts: []core.VolumeMount{
										core.VolumeMount{
											Name:      tempVolumeName,
											MountPath: "/app/temp_audio",
										},
									},
								},
							},
							Volumes: []core.Volume{
								core.Volume{
									Name: tempVolumeName,
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

	workerDeployment := utils.ManifestConfig{
		Filename: "worker-deployment.yaml",
		Manifests: []any{
			apps.NewDaemonSet(
				meta.ObjectMeta{
					Name: "audiomuse-ai-worker",
					Labels: map[string]string{
						"app.kubernetes.io/name":    fmt.Sprintf("%v-worker", generatorMeta.Name),
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				},
				apps.DaemonSetSpec{
					Selector: meta.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":    fmt.Sprintf("%v-worker", generatorMeta.Name),
							"app.kubernetes.io/version": generatorMeta.Docker.Version,
						},
					},
					Template: core.PodTemplateSpec{
						Metadata: meta.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":    fmt.Sprintf("%v-worker", generatorMeta.Name),
								"app.kubernetes.io/version": generatorMeta.Docker.Version,
							},
						},
						Spec: core.PodSpec{
							Containers: []core.Container{
								{
									Name:  "audiomuse-ai-worker",
									Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Env: []core.Env{
										{
											Name:  "SERVICE_TYPE",
											Value: "worker",
										},
										{
											Name:  "CLAP_ENABLED",
											Value: "false",
										},
										{
											Name: "JELLYFIN_USER_ID",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: secretName,
													Key:  "userid",
												},
											},
										},
										{
											Name: "JELLYFIN_TOKEN",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: secretName,
													Key:  "token",
												},
											},
										},
										{
											Name: "POSTGRES_USER",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: shared.PostgresCredsSecret,
													Key:  "username",
												},
											},
										},
										{
											Name: "POSTGRES_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: shared.PostgresCredsSecret,
													Key:  "password",
												},
											},
										},
									},
									EnvFrom: []core.EnvReference{
										{
											ConfigMapRef: core.ConfigMapRef{
												Name: configMapName,
											},
										},
									},
									Ports: []core.Port{
										{
											Name:          "worker",
											ContainerPort: generatorMeta.Port,
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      tempVolumeName,
											MountPath: "/app/temp_audio",
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: tempVolumeName,
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
						"app.kubernetes.io/name":    fmt.Sprintf("%v-flask", generatorMeta.Name),
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
					Ports: []core.ServicePort{
						{
							Name:       "http-audiomuse",
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
			flaskDeployment.Filename,
			workerDeployment.Filename,
			service.Filename,
			configMap.Filename,
			jellyfinCredsVaultSecret.Filename,
			pvc.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, workerDeployment, flaskDeployment, service, configMap, jellyfinCredsVaultSecret, pvc})
}
