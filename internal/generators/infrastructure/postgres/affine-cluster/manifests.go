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

func createAffineManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	cluster := utils.ManifestConfig{
		Filename: "affine-pg-cluster.yaml",
		Manifests: []any{
			cnpg.NewCluster(meta.ObjectMeta{
				Name: generatorMeta.Name,
			}, cnpg.ClusterSpec{
				Instances: 1,
				Storage: cnpg.ClusterStorage{
					StorageClass: shared.NFSLocalClass,
					Size:         "20Gi",
				},
				SuperuserSecret: cnpg.SuperuserSecret{
					Name: shared.AffinePGCredsSecret,
				},
				EnableSuperuserAccess: true,
			}),
		},
	}

	affineDB := utils.ManifestConfig{
		Filename: "affine-db.yaml",
		Manifests: []any{
			cnpg.NewDatabase(meta.ObjectMeta{
				Name: "affine-db",
			}, cnpg.DatabaseSpec{
				Name: "affine-db",
				Cluster: cnpg.DatabaseCluster{
					Name: generatorMeta.Name,
				},
				Owner: "postgres",
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
												"affine",
												"cloudnative-pg",
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
												"affine",
												"cnpg-system",
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
				affineDB.Filename,
				networkPolicy.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, cluster, affineDB, networkPolicy})
}
