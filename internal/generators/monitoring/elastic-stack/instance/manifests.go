package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"os"
)

func createElasticStackManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	vectorRoleConfig, err := os.ReadFile("./vector-role.yaml")
	if err != nil {
		fmt.Printf("Error while reading vector-role,yaml")
	}

	vectorRoleName := "vector-role"
	vectorRole := utils.ManifestConfig{
		Filename: "vector-role.yaml",
		Manifests: []any{
			core.NewSecret(meta.ObjectMeta{
				Name: vectorRoleName,
			}, core.SecretConfig{
				StringData: map[string]string{
					"roles.yml": string(vectorRoleConfig),
				},
			}),
		},
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"eck-elasticsearch": map[string]any{
				"enabled":          true,
				"fullnameOverride": "elasticsearch",
				"auth": map[string]any{
					"fileRealm": []map[string]string{
						{
							"secretName": generators.ElasticSearchAdminSecretName,
						},
						{
							"secretName": generators.ElasticSearchVectorSecretName,
						},
					},
					"roles": []map[string]string{
						{
							"secretName": vectorRoleName,
						},
					},
				},
				"nodeSets": []map[string]any{
					{
						"name":  "default",
						"count": 1,
						"config": map[string]any{
							"node.store.allow_mmap": false,
						},
						"volumeClaimTemplates": []map[string]any{
							{
								"metadata": map[string]any{
									"name": "elasticsearch-data",
								},
								"spec": map[string]any{
									"accessModes": []string{
										"ReadWriteOnce",
									},
									"storageClassName": generators.NFSRemoteClass,
									"resources": map[string]any{
										"requests": map[string]any{
											"storage": "20Gi",
										},
									},
								},
							},
						},
						"podTemplate": core.PodTemplateSpec{
							Spec: core.PodSpec{
								Containers: []core.Container{
									{
										Name: "elasticsearch",
										Resources: core.Resources{
											Limits: map[string]string{
												"memory": "2Gi",
											},
											Requests: map[string]string{
												"memory": "2Gi",
											},
										},
									},
								},
							},
						},
					},
				},
				"http": map[string]any{
					"tls": map[string]any{
						"selfSignedCertificate": map[string]any{
							"disabled": true,
						},
					},
				},
			},
			"eck-kibana": map[string]any{
				"enabled": true,
				"elasticsearchRef": map[string]any{
					"name": "elasticsearch",
				},
			},
		},
		nil,
	)

	service := utils.ManifestConfig{
		Filename: "service.yaml",
		Manifests: []any{
			core.NewService(
				meta.ObjectMeta{
					Name: "kibana",
					Labels: map[string]string{
						"app.kubernetes.io/name":    "eck-kibana",
						"app.kubernetes.io/version": generatorMeta.Helm.Version,
					},
					Annotations: map[string]string{
						"netbird.io/expose": "true",
						"netbird.io/groups": "cluster-services",
					},
				}, core.ServiceSpec{
					Selector: map[string]string{
						"common.k8s.elastic.co/type": "kibana",
					},
					Ports: []core.ServicePort{
						{
							Name:       "http-kibana",
							Port:       5601,
							TargetPort: 5601,
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
			repo.Filename,
			chart.Filename,
			release.Filename,
			vectorRole.Filename,
			service.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, vectorRole, service})
}
