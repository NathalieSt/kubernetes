package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/cnpg"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createPostgresManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	clusterName := "postgres"
	cluster := utils.ManifestConfig{
		Filename: "cnpg-cluster.yaml",
		Manifests: []any{
			cnpg.NewCluster(meta.ObjectMeta{
				Name: clusterName,
			}, cnpg.ClusterSpec{
				Instances: 3,
				//ImageName: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
				//Bootstrap: cnpg.Bootstrap{
				//	Initdb: cnpg.InitDB{
				//		PostInitTemplateSQL: []string{
				//			"CREATE EXTENSION postgis;",
				//			"CREATE EXTENSION postgis_topology;",
				//			"CREATE EXTENSION fuzzystrmatch;",
				//			"CREATE EXTENSION postgis_tiger_geocoder;",
				//		},
				//	},
				//},
				Storage: cnpg.ClusterStorage{
					StorageClass: generators.NFSRemoteClass,
					Size:         "15Gi",
				},
				InheritedMetadata: cnpg.InheritedMetadata{
					Annotations: map[string]string{
						"proxy.istio.io/config": "{\"holdApplicationUntilProxyStarts\": true}",
					},
				},
				SuperuserSecret: cnpg.SuperuserSecret{
					Name: generators.PostgresCredsSecret,
				},
				EnableSuperuserAccess: true,
			}),
		},
	}

	forgejoDB := utils.ManifestConfig{
		Filename: "forgejo-db.yaml",
		Manifests: []any{
			cnpg.NewDatabase(meta.ObjectMeta{
				Name: "forgejo-db",
			}, cnpg.DatabaseSpec{
				Name: "forgejo",
				Cluster: cnpg.DatabaseCluster{
					Name: clusterName,
				},
				Owner: "postgres",
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

	peerAuth := utils.ManifestConfig{
		Filename: "peer-auth.yaml",
		Manifests: []any{
			istio.NewPeerAuthenthication(meta.ObjectMeta{
				Name: "allow-permissive-postgres-access",
			}, istio.PeerAuthenthicationSpec{
				MTLS: istio.PeerAuthenthicationmTLS{
					Mode: istio.PERMISSIVE,
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
				forgejoDB.Filename,
				dawarichDB.Filename,
				peerAuth.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, cluster, forgejoDB, dawarichDB, peerAuth})
}
