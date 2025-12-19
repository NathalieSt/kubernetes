package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/cnpg"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createMatrixClusterManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	cluster := utils.ManifestConfig{
		Filename: "matrix-cluster.yaml",
		Manifests: []any{
			cnpg.NewCluster(meta.ObjectMeta{
				Name: generatorMeta.Name,
			}, cnpg.ClusterSpec{
				Instances: 1,
				Storage: cnpg.ClusterStorage{
					StorageClass: generators.DebianStorageClass,
					Size:         "100Gi",
				},
				Affinity: cnpg.AffinityConfiguration{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": "debian",
					},
				},
				SuperuserSecret: cnpg.SuperuserSecret{
					Name: generators.MatrixPGCredsSecret,
				},
				EnableSuperuserAccess: true,
			}),
		},
	}

	synapseDB := utils.ManifestConfig{
		Filename: "synapse-db.yaml",
		Manifests: []any{
			cnpg.NewDatabase(meta.ObjectMeta{
				Name: "synapse-db",
			}, cnpg.DatabaseSpec{
				Name: "synapse-db",
				Cluster: cnpg.DatabaseCluster{
					Name: generatorMeta.Name,
				},
				Owner: "postgres",
			}),
		},
	}

	discordBridgeDB := utils.ManifestConfig{
		Filename: "discord-bridge-db.yaml",
		Manifests: []any{
			cnpg.NewDatabase(meta.ObjectMeta{
				Name: "discord-bridge-db",
			}, cnpg.DatabaseSpec{
				Name: "discord-bridge-db",
				Cluster: cnpg.DatabaseCluster{
					Name: generatorMeta.Name,
				},
				Owner: "postgres",
			}),
		},
	}

	whatsAppBridgeDB := utils.ManifestConfig{
		Filename: "whatsapp-bridge-db.yaml",
		Manifests: []any{
			cnpg.NewDatabase(meta.ObjectMeta{
				Name: "whatsapp-bridge-db",
			}, cnpg.DatabaseSpec{
				Name: "whatsapp-bridge-db",
				Cluster: cnpg.DatabaseCluster{
					Name: generatorMeta.Name,
				},
				Owner: "postgres",
			}),
		},
	}

	signalBridgeDB := utils.ManifestConfig{
		Filename: "signal-bridge-db.yaml",
		Manifests: []any{
			cnpg.NewDatabase(meta.ObjectMeta{
				Name: "signal-bridge-db",
			}, cnpg.DatabaseSpec{
				Name: "signal-bridge-db",
				Cluster: cnpg.DatabaseCluster{
					Name: generatorMeta.Name,
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
				synapseDB.Filename,
				discordBridgeDB.Filename,
				whatsAppBridgeDB.Filename,
				signalBridgeDB.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{
		namespace,
		kustomization,
		cluster,
		synapseDB,
		discordBridgeDB,
		whatsAppBridgeDB,
		signalBridgeDB,
	})
}
