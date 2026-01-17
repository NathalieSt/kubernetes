package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
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
				StorageClassName: shared.NFSLocalClass,
			}),
		},
	}

	hwaVolumeName := "hwa-volume"
	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"image": map[string]any{
				"tag": "10.10.7",
			},
			"persistence": map[string]any{
				"media": map[string]any{
					"existingClaim": pvcName,
				},
				"config": map[string]any{
					"storageClass": shared.NFSLocalClass,
				},
			},
			"nodeSelector": map[string]any{
				"kubernetes.io/hostname": "debian",
			},
			"resources": core.Resources{
				Requests: map[string]string{
					"gpu.intel.com/i915": "1",
				},
				Limits: map[string]string{
					"gpu.intel.com/i915": "1",
				},
			},
			"securityContext": core.ContainerSecurityContext{
				Capabilities: core.Capabilities{
					Add: []string{"SYS_ADMIN"},
				},
				Privileged: false,
			},
			"extraVolumes": core.Volume{
				Name: hwaVolumeName,
				HostPath: core.HostPath{
					Path: "/dev/dri",
				},
			},
			"extraVolumeMounts": core.VolumeMount{
				Name:      hwaVolumeName,
				MountPath: "/dev/dri",
			},
		},
		nil,
	)

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
										"app.kubernetes.io/name": "flood",
									},
								},
							},
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

	jellyfinAudiomuseNetworkPolicy := utils.ManifestConfig{
		Filename: "jellyfin-audiomuse-network-policy.yaml",
		Manifests: []any{
			networking.NewNetworkPolicy(meta.ObjectMeta{
				Name: fmt.Sprintf("%v-audiomuse-networkpolicy", generatorMeta.Name),
			}, networking.NetworkPolicySpec{
				PolicyTypes: []networking.NetworkPolicyType{networking.Ingress},
				PodSelector: meta.LabelSelector{
					MatchLabels: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
				},
				Ingress: []networking.NetworkPolicyIngressRule{
					{
						From: []networking.NetworkPolicyPeer{
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "audiomuse-ai-flask",
									},
								},
							},
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "audiomuse-ai-worker",
									},
								},
							},
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "audiomuse-ai-nvidia-worker",
									},
								},
							},
						},
					},
				},
			}),
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
			scaledObject.Filename,
			networkPolicy.Filename,
			jellyfinAudiomuseNetworkPolicy.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, pvc, scaledObject, networkPolicy, jellyfinAudiomuseNetworkPolicy})
}
