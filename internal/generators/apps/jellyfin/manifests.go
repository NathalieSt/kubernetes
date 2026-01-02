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

func createJellyfinManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
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
			}),
		},
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"persistence": map[string]any{
				"media": map[string]any{
					"existingClaim": pvcName,
				},
				"config": map[string]any{
					"storageClass": generators.NFSLocalClass,
				},
			},
			"nodeSelector": map[string]any{
				"kubernetes.io/hostname": "debian",
			},
			"service": map[string]any{
				"annotations": map[string]any{
					"netbird.io/expose": "true",
					"netbird.io/groups": "cluster-services",
					//"netbird.io/resource-name": "jellyfin",
				},
			},
		},
		nil,
	)

	vpnSecretConfig := utils.StaticSecretConfig{
		Name:       fmt.Sprintf("%v-vpn", generatorMeta.Name),
		SecretName: fmt.Sprintf("%v-vpn", generatorMeta.Name),
		Path:       "vpn",
	}

	transConfigSecret := utils.StaticSecretConfig{
		Name:       "trans-config",
		SecretName: "trans-config",
		Path:       "vpn/vpn-nl-country-config",
	}

	vpnVaultSecret := utils.ManifestConfig{
		Filename: "vault-secret.yaml",
		Manifests: utils.GenerateVaultAccessManifests(
			generatorMeta.Name,
			//FIXME: get this from VSO generator meta
			"vault-secrets-operator",
			[]utils.StaticSecretConfig{vpnSecretConfig, transConfigSecret},
		),
	}

	transPVCVPNName := "trans-vpn-pvc"
	transPVCVPN := utils.ManifestConfig{
		Filename: "trans-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: transPVCVPNName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			}),
		},
	}

	transVPNVolume := "trans-vpn-volume"
	transVPNPVCVolume := "trans-vpn-vpc-volume"
	mediaVolume := "media-volume"
	deployment := utils.ManifestConfig{
		Filename: "deployment.yaml",
		Manifests: []any{
			apps.NewDeployment(
				meta.ObjectMeta{
					Name: "transsetup",
					Labels: map[string]string{
						"app.kubernetes.io/name":    "transsetup",
						"app.kubernetes.io/version": "1.0",
					},
				},
				apps.DeploymentSpec{
					Replicas: 1,
					Selector: meta.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":    "transsetup",
							"app.kubernetes.io/version": "1.0",
						},
					},
					Template: core.PodTemplateSpec{
						Metadata: meta.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":    "transsetup",
								"app.kubernetes.io/version": "1.0",
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
										`mkdir -p /tocopyto
										cp /readonly/* /tocopyto/
										chmod 644 /tocopyto/*
										chmod 755 /tocopyto/update-port.sh`,
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      transVPNVolume,
											MountPath: "/readonly",
										},
										{
											Name:      transVPNPVCVolume,
											MountPath: "/tocopyto",
										},
									},
								},
							},
							Containers: []core.Container{
								{
									Name:  "transmission-openvpn",
									Image: "jesec/flood:4.11",
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/data",
											Name:      mediaVolume,
										},
									},
									Ports: []core.Port{
										{
											Name:          "flood-web-ui",
											ContainerPort: 3000,
										},
									},
								},
								{
									Name:  "transmission-openvpn",
									Image: "haugene/transmission-openvpn:5.3.2",
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/data",
											Name:      mediaVolume,
										},
										{
											MountPath: "/etc/openvpn/custom/",
											Name:      transVPNPVCVolume,
										},
									},
									Ports: []core.Port{
										{
											Name:          "trans-web-ui",
											ContainerPort: 9091,
										},
									},
									Env: []core.Env{
										{
											Name:  "DEBUG",
											Value: "true",
										},
										{
											Name:  "OPENVPN_PROVIDER",
											Value: "custom",
										},
										{
											Name:  "OPENVPN_CONFIG",
											Value: "node-nl.protonvpn.udp",
										},
										{
											Name: "OPENVPN_USERNAME",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: vpnSecretConfig.SecretName,
													Key:  "user",
												},
											},
										},
										{
											Name: "OPENVPN_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: vpnSecretConfig.SecretName,
													Key:  "password",
												},
											},
										},
										{
											Name:  "LOCAL_NETWORK",
											Value: "10.0.0.0/24",
										},
									},
									SecurityContext: core.ContainerSecurityContext{
										Capabilities: core.Capabilities{
											Add: []string{
												"NET_ADMIN",
											},
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: transVPNPVCVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: transPVCVPNName,
									},
								},
								{
									Name: mediaVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: pvcName,
									},
								},
								{
									Name: transVPNVolume,
									Secret: core.SecretVolumeSource{
										SecretName: "trans-config",
										Items: []core.SecretVolumeItem{
											{
												Key:  "vpn-config",
												Path: "node-nl.protonvpn.udp.ovpn",
											},
											{
												Key:  "rotate-script",
												Path: "update-port.sh",
											},
										},
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
					Name: "transmission",
					Labels: map[string]string{
						"app.kubernetes.io/name":    "transsetup",
						"app.kubernetes.io/version": "1.0",
					},
				}, core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name":    "transsetup",
						"app.kubernetes.io/version": "1.0",
					},
					Ports: []core.ServicePort{
						{
							Name:       "http-trans-webui",
							Port:       9091,
							TargetPort: 9091,
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
			repo.Filename,
			chart.Filename,
			release.Filename,
			pvc.Filename,
			deployment.Filename,
			scaledObject.Filename,
			vpnVaultSecret.Filename,
			service.Filename,
			transPVCVPN.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, pvc, scaledObject, deployment, vpnVaultSecret, service, transPVCVPN})
}
