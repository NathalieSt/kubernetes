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

	livekitCMName := fmt.Sprintf("%v-config", generatorMeta.Name)
	livekitConfigMap := utils.ManifestConfig{
		Filename: "configmap.yaml",
		Manifests: []any{
			core.NewConfigMap(meta.ObjectMeta{
				Name: livekitCMName,
			}, map[string]string{
				"livekit.yaml": `
port: 7880
rtc:
  tcp_port: 7881
  port_range_start: 50000
  port_range_end: 50020
  use_external_ip: false

turn:
  enabled: true
  domain: livekit.nathi.lol
  tls_port: 5349
  udp_port: 3478
  external_tls: true

keys: {}`,
				"initialize.sh": `#!/bin/sh exec livekit-server --config /etc/livekit.yaml`,
			}),
		},
	}

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
									Command: []string{"livekit-server", "--config", "/etc/livekit.yaml"},
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
											Name:      livekitConfigVolumeName,
											MountPath: "/initialize.sh",
											SubPath:   "initialize.sh",
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: livekitConfigVolumeName,
									ConfigMap: core.ConfigMapVolumeSource{
										Name:        livekitCMName,
										DefaultMode: 0755,
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
				PolicyTypes: []networking.NetworkPolicyType{networking.Ingress},
				Ingress: []networking.NetworkPolicyIngressRule{
					{
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
							{
								// Allow all traffic within the same namespace
								PodSelector: meta.LabelSelector{},
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
			livekitConfigMap.Filename,
			livekitVaultSecret.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{
		namespace,
		kustomization,
		deployment,
		livekitConfigMap,
		service,
		networkPolicy,
		livekitVaultSecret,
	})
}
