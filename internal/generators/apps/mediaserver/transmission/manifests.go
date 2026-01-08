package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createTransmissionManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

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
				StorageClassName: shared.NFSLocalClass,
			}),
		},
	}

	transConfigPVCName := "trans-config-pvc"
	transConfigPVC := utils.ManifestConfig{
		Filename: "trans-config-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: transConfigPVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSLocalClass,
			}),
		},
	}

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

	ipLeakConfigMapName := "discord-bridge-configmap"
	ipLeakConfigMap, err := getIPLeakConfigMap(ipLeakConfigMapName)
	if err != nil {
		fmt.Println("An error occurred while getting the configMap for ip-leak")
		return nil
	}

	ipLeakConfigMapManifest := utils.ManifestConfig{
		Filename:  "configmap.yaml",
		Manifests: []any{*ipLeakConfigMap},
	}

	ipLeakVolume := "ip-leak-volume"
	transConfigVolume := "trans-config-volume"
	floodConfigVolume := "flood-config-volume"
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
							SecurityContext: core.PodSecurityContext{
								RunAsUser:  1000,
								RunAsGroup: 1001,
							},
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
								{
									Name:  "fix-file-permissions",
									Image: "alpine:latest",
									Command: []string{
										"/bin/sh",
										"-c",
										`chown -R 1000:1000 /config
										`,
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      floodConfigVolume,
											MountPath: "/config",
										},
									},
								},
							},
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
										{
											MountPath: "/config",
											Name:      transConfigVolume,
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
											Value: "10.244.0.0/16",
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
								{
									Name:    "ipleak",
									Image:   "curlimages/curl:8.17.0",
									Command: []string{"/bin/sh", "-c", "sh /scripts/ipleak.sh"},
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/scripts/",
											Name:      ipLeakVolume,
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: transConfigVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: transConfigPVCName,
									},
								},
								{
									Name: floodConfigVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: floodConfigPVCName,
									},
								},
								{
									Name: transVPNPVCVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: transPVCVPNName,
									},
								},
								{
									Name: mediaVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										//FIXME: get from jellyfin generator
										ClaimName: "jellyfin-pvc",
									},
								},
								{
									Name: ipLeakVolume,
									ConfigMap: core.ConfigMapVolumeSource{
										Name: ipLeakConfigMapName,
										Items: []core.VolumeConfigMapItem{
											{
												Key:  "ipleak.sh",
												Path: "ipleak.sh",
											},
										},
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
			vpnVaultSecret.Filename,
			service.Filename,
			transPVCVPN.Filename,
			transConfigPVC.Filename,
			floodConfigPVC.Filename,
			ipLeakConfigMapManifest.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, deployment, vpnVaultSecret, service, transPVCVPN, transConfigPVC, floodConfigPVC, ipLeakConfigMapManifest})
}
