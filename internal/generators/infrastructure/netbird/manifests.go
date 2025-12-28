package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createNetbirdManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

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
											Value: "netbird-k8s-router",
										},
										{
											Name:  "NB_LOG_LEVEL",
											Value: "debug",
										},
										{
											Name:  "ENABLE_ROSENPASS",
											Value: "true",
										},
									},
									SecurityContext: core.ContainerSecurityContext{
										Capabilities: core.Capabilities{
											Add: []string{
												"NET_ADMIN",
												"SYS_RESOURCE",
												"SYS_ADMIN",
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

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			deployment.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment})
}
