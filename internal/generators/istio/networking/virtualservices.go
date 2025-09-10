package main

import (
	"fmt"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func getVirtualServices(exposedServices []generator.GeneratorMeta, ingressGatewayName string) ([]any, error) {
	virtualServices := []any{}

	for _, service := range exposedServices {
		virtualServices = append(virtualServices, istio.NewVirtualService(
			meta.ObjectMeta{
				Name: fmt.Sprintf("%v-virtualservice", service.Name),
			},
			istio.VirtualServiceSpec{
				Hosts:    []string{fmt.Sprintf("%v.netbird.selfhosted", service.Caddy.DNSName)},
				Gateways: []string{ingressGatewayName},
				HTTP: []istio.VirtualServiceHTTP{
					{
						Route: []istio.VirtualServiceRoute{
							{
								Destination: istio.VirtualServiceDestination{
									Host: service.ClusterUrl,
									Port: istio.VirtualServicePortSelector{
										Number: service.Port,
									},
								},
								Headers: istio.VirtualServiceHeaders{
									Request: istio.VirtualServiceRequest{
										Set: istio.VirtualServiceSet{
											XForwardedProto: istio.VirtualServiceHTTPSProto,
										},
									},
								},
							},
						},
					},
				},
			},
		))
	}

	return virtualServices, nil
}
