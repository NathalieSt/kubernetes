package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/cnpg"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createSynapseClusterManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	cluster := utils.ManifestConfig{
		Filename: "synapse-cluster.yaml",
		Manifests: []any{
			cnpg.NewCluster(meta.ObjectMeta{
				Name: generatorMeta.Name,
			}, cnpg.ClusterSpec{
				Instances: 1,
				InheritedMetadata: cnpg.InheritedMetadata{
					Annotations: map[string]string{
						"proxy.istio.io/config": "{\"holdApplicationUntilProxyStarts\": true}",
					},
				},
				Storage: cnpg.ClusterStorage{
					StorageClass: generators.NFSLocalClass,
					Size:         "100Gi",
				},
				SuperuserSecret: cnpg.SuperuserSecret{
					Name: generators.SynapsePGCredsSecret,
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

	peerAuth := utils.ManifestConfig{
		Filename: "peer-auth.yaml",
		Manifests: []any{
			istio.NewPeerAuthenthication(meta.ObjectMeta{
				Name: "cnpg-permissive-mtls",
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
				synapseDB.Filename,
				peerAuth.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, cluster, synapseDB, peerAuth})
}
