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
)

func createDCTSManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	pvcAppSvName := "dcts-app-sv"
	pvcAppSv := utils.ManifestConfig{
		Filename: "pvc-app-sv.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcAppSvName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}

	pvcAppConfigsName := "dcts-app-configs"
	pvcAppConfigs := utils.ManifestConfig{
		Filename: "pvc-app-configs.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcAppConfigsName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}

	pvcAppUploadsName := "dcts-app-uploads"
	pvcAppUploads := utils.ManifestConfig{
		Filename: "pvc-app-uploads.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcAppUploadsName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "10Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}

	pvcAppEmojisName := "dcts-app-emojis"
	pvcAppEmojis := utils.ManifestConfig{
		Filename: "pvc-app-emojis.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcAppEmojisName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}

	pvcAppPluginsName := "dcts-app-plugins"
	pvcAppPlugins := utils.ManifestConfig{
		Filename: "pvc-app-plugins.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcAppPluginsName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}

	pvcAppThemesName := "dcts-app-themes"
	pvcAppThemes := utils.ManifestConfig{
		Filename: "pvc-app-themes.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcAppThemesName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}

	pvcAppGraphsName := "dcts-app-graphs"
	pvcAppGraphs := utils.ManifestConfig{
		Filename: "pvc-app-graphs.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcAppGraphsName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}

	pvcAppBackupName := "dcts-app-backup"
	pvcAppBackup := utils.ManifestConfig{
		Filename: "pvc-app-backup.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcAppBackupName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "5Gi",
				}},
				StorageClassName: shared.NFSRemoteClass,
			}),
		},
	}
	secretName := fmt.Sprintf("%v-secret", generatorMeta.Name)
	dctsSecretConfig := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-secret", generatorMeta.Name),
		SecretName: secretName,
		Path:       "dcts_db",
	}

	dctsVaultSecret := utils.ManifestConfig{
		Filename: "vault-secret.yaml",
		Manifests: utils.GenerateVaultAccessManifests(
			generatorMeta.Name,
			//FIXME: get this from VSO generator meta
			"vault-secrets-operator",
			[]utils.StaticSecretConfig{dctsSecretConfig},
		),
	}

	appSvVolume := "app-sv-volume"
	appConfigsVolume := "app-configs-volume"
	appUploadsVolume := "app-uploads-volume"
	appEmojisVolume := "app-emojis-volume"
	appPluginsVolume := "app-plugins-volume"
	appThemesVolume := "app-themes-volume"
	livekitConfigVolumeName := "livekit-volume"
	appGraphsVolume := "app-graphs-volume"
	appBackupVolume := "app-backup-volume"
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
							SecurityContext: core.PodSecurityContext{
								FsGroup: 1000,
							},
							Containers: []core.Container{
								{
									Name:  generatorMeta.Name,
									Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Env: []core.Env{
										{Name: "DB_HOST", Value: "remote-fs.netbird.nathalie-stiefsohn.eu"},
										{
											Name: "DB_USER",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: secretName,
													Key:  "username",
												},
											},
										},
										{
											Name: "DB_PASS",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: secretName,
													Key:  "password",
												},
											},
										},
										{Name: "DB_NAME", Value: "dcts"},
										{Name: "DEBUG", Value: "false"},
										{
											Name:  "LIVEKIT_URL",
											Value: "livekit.nathalie-stiefsohn.eu",
										},
										{
											Name: "LIVEKIT_API_KEY",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: "livekit-secret",
													Key:  "api-key",
												},
											},
										},
										{
											Name: "LIVEKIT_API_SECRET",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: "livekit-secret",
													Key:  "api-secret",
												},
											},
										},
										{Name: "REDIS_HOST", Value: "redis.redis.svc.cluster.local"},
										{Name: "REDIS_PORT", Value: "6379"},
									},
									Ports: []core.Port{
										{ContainerPort: generatorMeta.Port, Name: "http"},
									},
									VolumeMounts: []core.VolumeMount{
										{Name: appSvVolume, MountPath: "/app/sv"},
										{Name: appConfigsVolume, MountPath: "/app/configs"},
										{Name: appUploadsVolume, MountPath: "/app/public/uploads"},
										{Name: appEmojisVolume, MountPath: "/app/public/emojis"},
										{Name: appPluginsVolume, MountPath: "/app/plugins"},
										{Name: appThemesVolume, MountPath: "/app/themes"},
										{Name: appGraphsVolume, MountPath: "/app/public/graphs"},
										{Name: appBackupVolume, MountPath: "/app/backups"},
										{
											Name:      livekitConfigVolumeName,
											MountPath: "/app/livekit/livekit.yaml",
											SubPath:   "livekit.yaml",
										},
									},
									SecurityContext: core.ContainerSecurityContext{
										RunAsUser:                1000,
										RunAsNonRoot:             true,
										AllowPrivilegeEscalation: false,
										ReadOnlyRootFilesystem:   true,
										Capabilities: core.Capabilities{
											Drop: []string{"ALL"},
										},
									},
								},
							},
							Volumes: []core.Volume{
								{Name: appSvVolume, PersistentVolumeClaim: core.PVCVolumeSource{ClaimName: pvcAppSvName}},
								{Name: appConfigsVolume, PersistentVolumeClaim: core.PVCVolumeSource{ClaimName: pvcAppConfigsName}},
								{Name: appUploadsVolume, PersistentVolumeClaim: core.PVCVolumeSource{ClaimName: pvcAppUploadsName}},
								{Name: appEmojisVolume, PersistentVolumeClaim: core.PVCVolumeSource{ClaimName: pvcAppEmojisName}},
								{Name: appPluginsVolume, PersistentVolumeClaim: core.PVCVolumeSource{ClaimName: pvcAppPluginsName}},
								{Name: appThemesVolume, PersistentVolumeClaim: core.PVCVolumeSource{ClaimName: pvcAppThemesName}},
								{Name: appGraphsVolume, PersistentVolumeClaim: core.PVCVolumeSource{ClaimName: pvcAppGraphsName}},
								{Name: appBackupVolume, PersistentVolumeClaim: core.PVCVolumeSource{ClaimName: pvcAppBackupName}},
								{
									Name: livekitConfigVolumeName,
									Secret: core.SecretVolumeSource{
										SecretName:  "livekit-secret",
										DefaultMode: 0644,
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

	networkPolicy := utils.ManifestConfig{
		Filename: "network-policy.yaml",
		Manifests: []any{
			networking.NewNetworkPolicy(meta.ObjectMeta{
				Name: fmt.Sprintf("%v-networkpolicy", generatorMeta.Name),
			}, networking.NetworkPolicySpec{
				PolicyTypes: []networking.NetworkPolicyType{networking.Ingress, networking.Egress},
				PodSelector: meta.LabelSelector{
					MatchLabels: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
				},
				// Access from Netbird
				Ingress: []networking.NetworkPolicyIngressRule{
					{
						From: []networking.NetworkPolicyPeer{
							shared.ExternalIngressNetworkPolicyPeer,
						},
						Ports: []networking.NetworkPolicyPort{
							{
								Port:     int32(generatorMeta.Port),
								Protocol: networking.TCP,
							},
						},
					},
				},
				Egress: []networking.NetworkPolicyEgressRule{
					// Access to CoreDNS
					utils.GetCoreDNSEgressRule(),
					// Access to MariaDB
					{
						To: []networking.NetworkPolicyPeer{
							{
								IpBlock: networking.IPBlock{
									CIDR: "100.87.188.103/32",
								},
							},
						},
					},
					{
						Ports: []networking.NetworkPolicyPort{
							{
								Port:     3306,
								Protocol: networking.TCP,
							},
						},
					},
					// Access to Livekit
					{
						To: []networking.NetworkPolicyPeer{
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "livekit",
									},
								},
							},
						},
						Ports: []networking.NetworkPolicyPort{
							{
								Port:     7880,
								Protocol: networking.TCP,
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
			deployment.Filename,
			service.Filename,
			pvcAppConfigs.Filename,
			pvcAppEmojis.Filename,
			pvcAppPlugins.Filename,
			pvcAppSv.Filename,
			pvcAppThemes.Filename,
			pvcAppUploads.Filename,
			dctsVaultSecret.Filename,
			networkPolicy.Filename,
			pvcAppGraphs.Filename,
			pvcAppBackup.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, deployment, service, pvcAppConfigs, pvcAppBackup, pvcAppGraphs, pvcAppEmojis, pvcAppEmojis, pvcAppPlugins, pvcAppSv, pvcAppThemes, pvcAppUploads, dctsVaultSecret, networkPolicy})
}
