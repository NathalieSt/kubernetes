package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createVaultManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"ui": map[string]any{
				"enabled": true,
			},
			"server": map[string]any{
				"dataStorage": map[string]any{
					"storageClass": generators.NFSLocalClass,
				},
			},
		},
		nil,
	)

	peerAuth := utils.ManifestConfig{
		Filename: "peer-auth.yaml",
		Manifests: []any{
			istio.NewPeerAuthenthication(meta.ObjectMeta{
				Name: "vault-permissive-mtls",
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
				repo.Filename,
				chart.Filename,
				release.Filename,
				peerAuth.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, peerAuth})
}
