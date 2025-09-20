package main

import (
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createCNPGManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, false),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm, nil, nil)

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
				repo.Filename,
				chart.Filename,
				release.Filename,
				peerAuth.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, peerAuth})
}
