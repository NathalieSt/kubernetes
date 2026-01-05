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
	"kubernetes/pkg/schema/k8s/networking"
)

func createStashManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
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
					"storage": "200Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			}),
		},
	}

	volume := "volume"
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
											ContainerPort: generatorMeta.Port,
											Name:          generatorMeta.Name,
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      volume,
											MountPath: "/root/.stash",
											SubPath:   "config",
										},
										{
											Name:      volume,
											MountPath: "/data",
											SubPath:   "data",
										},
										{
											Name:      volume,
											MountPath: "/metadata",
											SubPath:   "metadata",
										},
										{
											Name:      volume,
											MountPath: "/blobs",
											SubPath:   "blobs",
										},
										{
											Name:      volume,
											MountPath: "/generated",
											SubPath:   "generated",
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: volume,
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

	scaledObject := utils.ManifestConfig{
		Filename:  "scaled-object.yaml",
		Manifests: utils.GenerateCronScaler(fmt.Sprintf("%v-scaledobject", generatorMeta.Name), generatorMeta.Name, keda.Deployment, generatorMeta.KedaScaling),
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
			deployment.Filename,
			service.Filename,
			scaledObject.Filename,
			networkPolicy.Filename,
			pvc.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment, service, scaledObject, networkPolicy, pvc})
}
