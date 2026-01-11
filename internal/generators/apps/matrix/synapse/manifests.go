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

func createSynapseManifests(generatorMeta generator.GeneratorMeta, rootDir string, relativeDir string) (map[string][]byte, error) {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	configMapName := "synapse-configmap"
	configMap, err := getSynapseConfigMap(configMapName, rootDir, relativeDir)
	if err != nil {
		fmt.Println("An error occurred while getting the configMap for synapse")
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
					"storage": "100Gi",
				}},
				StorageClassName: shared.DebianStorageClass,
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
	secretVolumeName := "signing-key-volume"
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
							NodeSelector: map[string]string{
								"kubernetes.io/hostname": "debian",
							},
							SecurityContext: core.PodSecurityContext{
								FsGroup: 991,
							},
							InitContainers: []core.Container{
								{
									Name:  "config-init",
									Image: "alpine:latest",
									Command: []string{
										"/bin/sh",
										"-c",
										`apk update && apk add gettext;
envsubst < /template/homeserver.yaml > /data/homeserver.yaml;
envsubst < /template/discord-registration.yaml > /data/discord-registration.yaml;
cp /template/matrix.cluster.netbird.selfhosted.log.config /data;
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
													Name: shared.DiscordBridgeSecretName,
												},
											},
										},
										{
											Name: "HS_TOKEN",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "hs_token",
													Name: shared.DiscordBridgeSecretName,
												},
											},
										},
										{
											Name: "SENDER_LOCALPART",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "sender_localpart",
													Name: shared.DiscordBridgeSecretName,
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
											Value: "synapse-db",
										},
										{
											Name: "REGISTRATION_SHARED_SECRET",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: shared.SynapseSecretName,
													Key:  "registration_shared_secret",
												},
											},
										},
										{
											Name: "MACAROON_SECRET_KEY",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: shared.SynapseSecretName,
													Key:  "macaroon_secret_key",
												},
											},
										},
										{
											Name: "FORM_SECRET",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: shared.SynapseSecretName,
													Key:  "form_secret",
												},
											},
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
										{
											Name:      secretVolumeName,
											MountPath: "/signing",
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
								{
									Name: secretVolumeName,
									Secret: core.SecretVolumeSource{
										SecretName: shared.SynapseSecretName,
										Items: []core.SecretVolumeItem{
											{
												Key:  "signing-key",
												Path: "matrix.cluster.netbird.selfhosted.signing.key",
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
			datapvc.Filename,
			service.Filename,
			configMapManifest.Filename,
			networkPolicy.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment, service, configMapManifest, datapvc, networkPolicy}), nil
}
