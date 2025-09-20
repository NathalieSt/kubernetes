package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"path"
)

func createSynapseManifests(generatorMeta generator.GeneratorMeta, rootDir string) (map[string][]byte, error) {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	configMapName := "synapse-configmap"
	configMap := utils.ManifestConfig{
		Filename:  "configmap.yaml",
		Manifests: []any{getSynapseConfigMap(configMapName)},
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
			},
			),
		},
	}

	datapvcName := fmt.Sprintf("%v-data-pvc", generatorMeta.Name)
	datapvc := utils.ManifestConfig{
		Filename: "data-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: generators.NFSLocalClass,
			},
			),
		},
	}

	postgresMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/postgres/synapse-cluster"))
	if err != nil {
		fmt.Println("An error happened while getting postgres meta ")
		return nil, err
	}

	configVolumeName := "config-volume"
	dataVolumeName := "data-volume"
	volumeName := "synapse-pvc-volume"
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
							InitContainers: []core.Container{
								{
									Name:    "config-init",
									Image:   "alpine:latest",
									Command: []string{"/bin/sh", "-c"},
									Args: []string{
										"envsubst < /template/homeserver.yaml > /data/homeserver.yaml; cp /template/matrix.cluster.netbird.selfhosted.log.config /data",
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
									Env: []core.Env{
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
													Name: generators.SynapsePGCredsSecret,
												},
											},
										},
										{
											Name: "POSTGRES_PASSWORD",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Key:  "password",
													Name: generators.SynapsePGCredsSecret,
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
													Name: generators.SynapseSecretName,
													Key:  "registration_shared_secret",
												},
											},
										},
										{
											Name: "MACAROON_SECRET_KEY",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: generators.SynapseSecretName,
													Key:  "macaroon_secret_key",
												},
											},
										},
										{
											Name: "FORM_SECRET",
											ValueFrom: core.ValueFrom{
												SecretKeyRef: core.SecretKeyRef{
													Name: generators.SynapseSecretName,
													Key:  "form_secret",
												},
											},
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      volumeName,
											MountPath: "/media",
										},
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
									Name: volumeName,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: pvcName,
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
										SecretName: generators.SynapseSecretName,
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

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			deployment.Filename,
			pvc.Filename,
			datapvc.Filename,
			service.Filename,
			configMap.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment, pvc, service, configMap, datapvc}), nil
}
