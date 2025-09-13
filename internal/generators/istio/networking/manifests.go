package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createIstioNetworkingManifests(rootDir string, generatorMeta generator.GeneratorMeta) (map[string][]byte, error) {

	exposedServices, err := utils.GetMetaForExposedGenerators(rootDir)
	if err != nil {
		fmt.Println("An error happened while getting exposed services")
		return nil, err
	}

	hosts := []string{}
	for _, service := range exposedServices {
		hosts = append(hosts, fmt.Sprintf("%v.netbird.selfhosted", service.Caddy.DNSName))
	}

	ingressClusterGatewayName := "ingress-cluster-gateway"
	ingressClusterGateway := utils.ManifestConfig{
		Filename: "ingress-cluster-gateway.yaml",
		Manifests: []any{
			istio.NewGateway(
				meta.ObjectMeta{
					Name: ingressClusterGatewayName,
				},
				istio.GatewaySpec{
					Servers: []istio.GatewayServer{
						{
							Port: istio.IstioPort{
								Number:   80,
								Name:     "http",
								Protocol: istio.HTTP,
							},
							Hosts: hosts,
						},
					},
					Selector: istio.GatewaySelector{
						Istio: "ingress",
					},
				},
			),
		},
	}

	virtualServiceManifests, err := getVirtualServices(exposedServices, ingressClusterGatewayName)
	if err != nil {
		fmt.Println("An error happened while getting VirtualServices")
		return nil, err
	}

	virtualServices := utils.ManifestConfig{
		Filename:  "virtual-services.yaml",
		Manifests: virtualServiceManifests,
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			virtualServices.Filename,
			ingressClusterGateway.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, virtualServices, ingressClusterGateway}), nil
}
