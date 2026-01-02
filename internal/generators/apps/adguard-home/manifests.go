package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
)

func createBookloreManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
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
				StorageClassName: generators.NFSRemoteClass,
			},
			),
		},
	}

	workPVCName := "work-pvc"
	workPVC := utils.ManifestConfig{
		Filename: "work-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: workPVCName,
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

	confVolume := "conf-volume"
	workVolume := "work-volume"
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
											ContainerPort: 54,
											Name:          "dns",
										},
										{
											ContainerPort: 3000,
											Name:          "setup-ui",
										},
										{
											ContainerPort: 80,
											Name:          "web-ui",
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/opt/adguardhome/conf",
											Name:      confVolume,
										},
										{
											MountPath: "/opt/adguardhome/work",
											Name:      workVolume,
										},
									},
								},
								{
									Name:  "netbird-agent",
									Image: "netbirdio/netbird:latest",
									Env: []core.Env{
										{
											Name: "NB_SETUP_KEY",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: generators.NetbirdSecretName,
													Key:  "setup-key",
												},
											},
										},
										{
											Name:  "NB_HOSTNAME",
											Value: "adguard-home",
										},
										{
											Name:  "NB_MANAGEMENT_URL",
											Value: "https://netbird.nathalie-stiefsohn.eu",
										},
									},
									Resources: core.Resources{
										Requests: map[string]string{
											"cpu":    "50m",
											"memory": "64Mi",
										},
										Limits: map[string]string{
											"cpu":    "100m",
											"memory": "128Mi",
										},
									},
									SecurityContext: core.ContainerSecurityContext{
										Privileged: true,
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
									Name: workVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: workPVCName,
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
								IpBlock: networking.IPBlock{
									CIDR: "100.127.0.0/16",
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
			workPVC.Filename,
			service.Filename,
			deployment.Filename,
			networkPolicy.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, workPVC, confPVC, kustomization, deployment, service, networkPolicy})
}
