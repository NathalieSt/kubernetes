package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
)

func createLivekitManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	secretName := fmt.Sprintf("%v-secret", generatorMeta.Name)
	livekitVaultSecret := utils.ManifestConfig{
		Filename: "vault-secret.yaml",
		Manifests: utils.GenerateVaultAccessManifests(
			generatorMeta.Name,
			"vault-secrets-operator",
			[]utils.StaticSecretConfig{
				{
					Name:       fmt.Sprintf("%v-secret", generatorMeta.Name),
					SecretName: secretName,
					Path:       "livekit",
				},
			},
		),
	}

	livekitCertVolumeName := "livekit-certs-volume"
	livekitConfigVolumeName := "livekit-config-volume"
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
									Name:    generatorMeta.Name,
									Image:   fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Command: []string{"/livekit-server", "--config", "/etc/livekit.yaml", "--port", "7880"},
									Ports: []core.Port{
										{
											ContainerPort: generatorMeta.Port,
											Name:          "http",
										},
									},
									Env: []core.Env{
										{
											Name: "LIVEKIT_API_KEY",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: secretName,
													Key:  "api-key",
												},
											},
										},
										{
											Name: "LIVEKIT_API_SECRET",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: secretName,
													Key:  "api-secret",
												},
											},
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      livekitConfigVolumeName,
											MountPath: "/etc/livekit.yaml",
											SubPath:   "livekit.yaml",
										},
										{
											Name:      livekitCertVolumeName,
											MountPath: "/certs",
											Readonly:  true,
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
								{
									Name: livekitConfigVolumeName,
									Secret: core.SecretVolumeSource{
										SecretName:  secretName,
										DefaultMode: 0644,
									},
								},
								{
									Name: livekitCertVolumeName,
									Secret: core.SecretVolumeSource{
										SecretName: "livekit-tls",
									},
								},
							},
						},
					},
				},
			),
		},
	}

	ports := []core.ServicePort{
		{
			Name:       fmt.Sprintf("tcp-%v", generatorMeta.Name),
			Port:       7880,
			TargetPort: 7880,
			Protocol:   core.TCP,
		},
		{
			Name:       fmt.Sprintf("tcp-rtc-%v", generatorMeta.Name),
			Port:       7881,
			TargetPort: 7881,
			Protocol:   core.TCP,
		},
		{
			Name:       fmt.Sprintf("udp-turn-%v", generatorMeta.Name),
			Port:       3478,
			TargetPort: 3478,
			Protocol:   core.UDP,
		},
		{
			Name:       fmt.Sprintf("tcp-turns-%v", generatorMeta.Name),
			Port:       5349,
			TargetPort: 5349,
			Protocol:   core.TCP,
		},
		{
			Name:       fmt.Sprintf("udp-turns-%v", generatorMeta.Name),
			Port:       5349,
			TargetPort: 5349,
			Protocol:   core.UDP,
		},
	}

	for i := int64(50000); i <= 50004; i++ {
		ports = append(ports, core.ServicePort{
			Name:       fmt.Sprintf("rtc-udp-%d", i),
			Port:       i,
			TargetPort: i,
			Protocol:   core.UDP,
		})
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
					Ports: ports,
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
				PodSelector: meta.LabelSelector{
					MatchLabels: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
				},
				PolicyTypes: []networking.NetworkPolicyType{networking.Ingress, networking.Egress},
				Ingress: []networking.NetworkPolicyIngressRule{
					{
						// Access from Netbird
						From: []networking.NetworkPolicyPeer{
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "netbird-router",
									},
								},
								NamespaceSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"kubernetes.io/metadata.name": "netbird",
									},
								},
							},
						},
						Ports: []networking.NetworkPolicyPort{
							{
								Port:     7880,
								Protocol: networking.TCP,
							},
							{
								Port:     7881,
								Protocol: networking.TCP,
							},
						},
					},
					{
						// Access from DCTS
						From: []networking.NetworkPolicyPeer{
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "dcts",
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
					// Access from Clients for RTC UDP
					{
						Ports: []networking.NetworkPolicyPort{
							{
								Port:    50000,
								EndPort: 50005,
							},
						},
					},
					{
						// TURN/TURNS from anywhere (clients are external)
						From: []networking.NetworkPolicyPeer{},
						Ports: []networking.NetworkPolicyPort{
							{
								Port:     3478,
								Protocol: networking.UDP,
							},
							{
								Port:     5349,
								Protocol: networking.TCP,
							},
							{
								Port:     5349,
								Protocol: networking.UDP,
							},
						},
					},
				},
				Egress: []networking.NetworkPolicyEgressRule{
					// Access to CoreDNS
					utils.GetCoreDNSEgressRule(),
					{
						// Access to Clients for UDP
						Ports: []networking.NetworkPolicyPort{
							{
								Port:    50000,
								EndPort: 50005,
							},
						},
					},
					{
						// TURN/TURNS from anywhere (clients are external)
						To: []networking.NetworkPolicyPeer{},
						Ports: []networking.NetworkPolicyPort{
							{
								Port:     3478,
								Protocol: networking.UDP,
							},
							{
								Port:     5349,
								Protocol: networking.TCP,
							},
							{
								Port:     5349,
								Protocol: networking.UDP,
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
			deployment.Filename,
			service.Filename,
			networkPolicy.Filename,
			livekitVaultSecret.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{
		namespace,
		kustomization,
		deployment,
		service,
		networkPolicy,
		livekitVaultSecret,
	})
}
