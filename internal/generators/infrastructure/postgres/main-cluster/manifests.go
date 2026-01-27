package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/cnpg"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
)

func createPostgresManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	clusterName := "postgres"
	cluster := utils.ManifestConfig{
		Filename: "cnpg-cluster.yaml",
		Manifests: []any{
			cnpg.NewCluster(meta.ObjectMeta{
				Name: clusterName,
			}, cnpg.ClusterSpec{
				Instances: 1,
				Storage: cnpg.ClusterStorage{
					StorageClass: shared.NFSRemoteClass,
					Size:         "15Gi",
				},
				SuperuserSecret: cnpg.SuperuserSecret{
					Name: shared.PostgresCredsSecret,
				},
				EnableSuperuserAccess: true,
				Probes: cnpg.ProbesConfiguration{
					Readiness: cnpg.ProbeWithStrategy{
						InitialDelaySeconds: 30,
						PeriodSeconds:       30,
						TimeoutSeconds:      30,
						FailureThreshold:    10,
					},
				},
			}),
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
									MachExpressions: []meta.MatchExpression{
										{
											Key:      "app.kubernetes.io/name",
											Operator: meta.In,
											Values: []string{
												"mealie",
												"audiomuse-ai-worker",
												"audiomuse-ai-nvidia-worker",
												"audiomuse-ai-flask",
												"cloudnative-pg",
												"openclarity",
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
												"mealie",
												"jellyfin",
												"cnpg-system",
												"openclarity",
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
				cluster.Filename,
				networkPolicy.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, cluster, networkPolicy})
}
