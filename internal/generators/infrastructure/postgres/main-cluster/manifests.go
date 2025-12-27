package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/cnpg"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
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
				Instances: 3,
				Storage: cnpg.ClusterStorage{
					StorageClass: generators.NFSRemoteClass,
					Size:         "15Gi",
				},
				SuperuserSecret: cnpg.SuperuserSecret{
					Name: generators.PostgresCredsSecret,
				},
				EnableSuperuserAccess: true,
				Probes: cnpg.ProbesConfiguration{
					Readiness: cnpg.ProbeWithStrategy{
						InitialDelaySeconds: 30,
						PeriodSeconds:       10,
						TimeoutSeconds:      5,
						FailureThreshold:    6,
					},
				},
			}),
		},
	}

	dawarichDB := utils.ManifestConfig{
		Filename: "dawarich-db.yaml",
		Manifests: []any{
			cnpg.NewDatabase(meta.ObjectMeta{
				Name: "dawarich-db",
			}, cnpg.DatabaseSpec{
				Name: "dawarich-development",
				Cluster: cnpg.DatabaseCluster{
					Name: clusterName,
				},
				Owner: "postgres",
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
				dawarichDB.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, cluster, dawarichDB})
}
