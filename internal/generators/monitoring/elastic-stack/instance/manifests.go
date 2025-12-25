package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
)

func createElasticStackManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"eck-elasticsearch": map[string]any{
				"enabled":          true,
				"fullnameOverride": "elasticsearch",
				"auth:": map[string]any{
					"fileRealm": []map[string]string{
						{
							"secretName": generators.ElasticSearchAdminSecretName,
						},
						{
							"secretName": generators.ElasticSearchVectorSecretName,
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

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
			chart.Filename,
			release.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release})
}
