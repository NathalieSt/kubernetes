package main

import (
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createSecurityManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	peerAuthenthication := utils.ManifestConfig{
		Filename: "peer-authenthication.yaml",
		Manifests: []any{
			istio.NewPeerAuthenthication(meta.ObjectMeta{
				Name: "enforce-mtls-in-mesh",
			}, istio.PeerAuthenthicationSpec{
				MTLS: istio.PeerAuthenthicationmTLS{
					Mode: istio.STRICT,
				},
			}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename:  "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{peerAuthenthication.Filename}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, peerAuthenthication})
}
