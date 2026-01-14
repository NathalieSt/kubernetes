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
	"path"
)

func createWhatsappBridgeManifests(generatorMeta generator.GeneratorMeta, rootDir string, relativeDir string) (map[string][]byte, error) {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	configMapName := "whatsapp-bridge-configmap"
	configMap, err := getWhatsappBridgeConfigMap(configMapName, rootDir, relativeDir)
	if err != nil {
		fmt.Println("An error occurred while getting the configMap for whatsapp-bridge")
		return nil, err
	}

	configMapManifest := utils.ManifestConfig{
		Filename:  "configmap.yaml",
		Manifests: []any{*configMap},
	}

	datapvcName := fmt.Sprintf("%v-data-pvc", generatorMeta.Name)
	datapvc := utils.ManifestConfig{
		Filename: "data-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: datapvcName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSLocalClass,
			},
			),
		},
	}

	postgresMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/postgres/matrix-cluster"))
	if err != nil {
		fmt.Println("An error happened while getting postgres meta ")
		return nil, err
	}

	configVolumeName := "config-volume"
	dataVolumeName := "data-volume"
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
							InitContainers: []core.Container{
								{
									Name:  "config-init",
									Image: "alpine:latest",
									Command: []string{
										"/bin/sh",
										"-c",
										`apk update && apk add gettext;
envsubst < /template/config.yaml > /data/config.yaml;
										`,
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      configVolumeName,
											MountPath: "/template",
										},
										{
											Name:      dataVolumeName,
											MountPath: "/data",
										},
									},
									Env: []core.Env{
										{
											Name: "AS_TOKEN",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "as_token",
													Name: shared.WhatsappBridgeSecretName,
												},
											},
										},
										{
											Name: "HS_TOKEN",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "hs_token",
													Name: shared.WhatsappBridgeSecretName,
												},
											},
										},
										{
											Name:  "POSTGRES_SERVER",
											Value: postgresMeta.ClusterUrl,
										},
										{
											Name:  "POSTGRES_PORT",
											Value: fmt.Sprintf("%v", postgresMeta.Port),
										},
										{
											Name: "POSTGRES_USERNAME",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "username",
													Name: shared.MatrixPGCredsSecret,
												},
											},
										},
										{
											Name: "POSTGRES_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "password",
													Name: shared.MatrixPGCredsSecret,
												},
											},
										},
										{
											Name:  "POSTGRES_DB",
											Value: "whatsapp-bridge-db",
										},
									},
								},
							},
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
											Name:      dataVolumeName,
											MountPath: "/data",
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: configVolumeName,
									ConfigMap: core.ConfigMapVolumeSource{
										Name: configMapName,
									},
								},
								{
									Name: dataVolumeName,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: datapvcName,
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
				PolicyTypes: []networking.NetworkPolicyType{networking.Ingress, networking.Egress},
				Ingress: []networking.NetworkPolicyIngressRule{
					{
						From: []networking.NetworkPolicyPeer{
							{
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "synapse",
									},
								},
								NamespaceSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"kubernetes.io/metadata.name": "synapse",
									},
								},
							},
						},
					},
				},
				//FIXME: get label from corresponding generators
				//Might need an update of GeneratorMeta struct with those labels
				Egress: []networking.NetworkPolicyEgressRule{
					utils.GetCoreDNSEgressRule(),
					{
						To: []networking.NetworkPolicyPeer{
							{
								NamespaceSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"kubernetes.io/metadata.name": "synapse",
									},
								},
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "synapse",
									},
								},
							},
						},
					},
					{
						To: []networking.NetworkPolicyPeer{
							{
								NamespaceSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"kubernetes.io/metadata.name": "matrix-pg-cluster",
									},
								},
								PodSelector: meta.LabelSelector{
									MatchLabels: map[string]string{
										"cnpg.io/cluster": "matrix-pg",
									},
								},
							},
						},
					},
					{
						To: []networking.NetworkPolicyPeer{
							{
								IpBlock: networking.IPBlock{
									CIDR: "0.0.0.0/0",
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
			datapvc.Filename,
			service.Filename,
			configMapManifest.Filename,
			networkPolicy.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment, service, configMapManifest, datapvc, networkPolicy}), nil
}
