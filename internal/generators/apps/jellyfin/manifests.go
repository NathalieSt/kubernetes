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

	quiConfigPVCName := fmt.Sprintf("%v-qui-config-pvc", generatorMeta.Name)
	quiConfigPVC := utils.ManifestConfig{
		Filename: "quic-config-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: quiConfigPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			}),
		},
	}

	qbitConfigPVCName := fmt.Sprintf("%v-qbit-config-pvc", generatorMeta.Name)
	qbitConfigPVC := utils.ManifestConfig{
		Filename: "qbit-config-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: qbitConfigPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
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

	vpnVaultSecret := utils.ManifestConfig{
		Filename: "vault-secret.yaml",
		Manifests: utils.GenerateVaultAccessManifests(
			generatorMeta.Name,
			//FIXME: get this from VSO generator meta
			"vault-secrets-operator",
			[]utils.StaticSecretConfig{vpnSecretConfig},
		),
	}

	qbitConfigVolume := "qbit-config-volume"
	quiConfigVolume := "qui-config-volume"
	mediaVolume := "media-volume"
	deployment := utils.ManifestConfig{
		Filename: "deployment.yaml",
		Manifests: []any{
			apps.NewDeployment(
				meta.ObjectMeta{
					Name: "qbitsetup",
					Labels: map[string]string{
						"app.kubernetes.io/name":    "qbitsetup",
						"app.kubernetes.io/version": "1.0",
					},
				},
				apps.DeploymentSpec{
					Replicas: 1,
					Selector: meta.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":    "qbitsetup",
							"app.kubernetes.io/version": "1.0",
						},
					},
					Template: core.PodTemplateSpec{
						Metadata: meta.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":    "qbitsetup",
								"app.kubernetes.io/version": "1.0",
							},
						},
						Spec: core.PodSpec{
							Containers: []core.Container{
								{
									Name:  "qui",
									Image: "ghcr.io/autobrr/qui:latest",
									Ports: []core.Port{
										core.Port{
											Name:          "qui",
											ContainerPort: 7476,
										},
									},
									VolumeMounts: []core.VolumeMount{
										core.VolumeMount{
											Name:      quiConfigVolume,
											MountPath: "/config",
										},
									},
								},
								{
									Name:  "qbittorrent",
									Image: "linuxserver/qbittorrent:5.1.4",
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/downloads",
											Name:      mediaVolume,
										},
										{
											MountPath: "/config",
											Name:      qbitConfigVolume,
										},
									},
									Ports: []core.Port{
										{
											Name:          "web-ui",
											ContainerPort: 8080,
										},
									},
								},
								{
									Name:  "glue-tun",
									Image: "qmcgaw/gluetun:v3.41",
									SecurityContext: core.ContainerSecurityContext{
										Capabilities: core.Capabilities{
											Add: []string{
												"NET_ADMIN",
											},
										},
									},
									Env: []core.Env{
										{
											Name:  "VPN_SERVICE_PROVIDER",
											Value: "protonvpn",
										},
										{
											Name:  "VPN_TYPE",
											Value: "openvpn",
										},
										{
											Name:  "SERVER_COUNTRIES",
											Value: "Netherlands",
										},
										{
											Name:  "VPN_PORT_FORWARDING",
											Value: "on",
										},
										{
											Name:  "VPN_PORT_FORWARDING_PROVIDER",
											Value: "protonvpn",
										},
										{
											Name: "OPENVPN_USER",
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
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: mediaVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: pvcName,
									},
								},
								{
									Name: quiConfigVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: quiConfigPVCName,
									},
								},
								{
									Name: qbitConfigVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: qbitConfigPVCName,
									},
								},
							},
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
			quiConfigPVC.Filename,
			qbitConfigPVC.Filename,
			vpnVaultSecret.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, pvc, scaledObject, deployment, quiConfigPVC, qbitConfigPVC, vpnVaultSecret})
}
