package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
)

func createVaultManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"ui": map[string]any{
				"enabled": true,
				"annotations": map[string]string{
					"netbird.io/expose": "true",
					"netbird.io/groups": "cluster-services",
				},
			},
			"server": map[string]any{
				"dataStorage": map[string]any{
					"storageClass": shared.NFSLocalClass,
				},
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
									MachExpressions: []meta.MatchExpression{
										{
											Key:      "app.kubernetes.io/name",
											Operator: meta.In,
											Values: []string{
												"caddy",
												"vault-secrets-operator",
											},
										},
									},
								},
								NamespaceSelector: meta.LabelSelector{
									MachExpressions: []meta.MatchExpression{
										{
											Key:      "kubernetes.io/metadata.name",
											Operator: meta.In,
											Values: []string{
												"caddy",
												"vault-secrets-operator",
											},
										},
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
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				namespace.Filename,
				repo.Filename,
				chart.Filename,
				release.Filename,
				networkPolicy.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, networkPolicy})
}
