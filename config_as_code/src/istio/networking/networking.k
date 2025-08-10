import ....lib.istio
import ....lib.k8s_wrapper
import ...shared
import yaml

schema VirtualServiceAndFile:
    file_name: str
    virtual_service: istio.VirtualService

out_dir = shared.cluster_path + "istio/networking/"
ingress_out_dir = out_dir + "ingress/"
ingress_virtual_service_out_dir = ingress_out_dir + "virtual-services/"
ingress_virtual_service_kustomization_path = ingress_out_dir + ingress_virtual_service_out_dir
# Filenames
kustomization_file = "kustomization.yaml"
ingress_cluster_gateway_file = "ingress-cluster-gateway.yaml"
# Configs
service_name = "istio-networking"
ingress_cluster_gateway_name = "ingress-cluster-gateway"
# k8s configuration
ingress_virtual_services: [VirtualServiceAndFile] = [{
    file_name: service.name + "-virtualservice.yaml"
    virtual_service: {
        metadata: {
            name: service.name + "-virtualservice"
        }
        spec: {
            hosts: [
                "${service.gateway_config.dns_name}.netbird.selfhosted"
            ]
            gateways: [
                ingress_cluster_gateway_name
            ]
            http: [
                {
                    destination: [
                        {
                            host: "${service.cluster_url}"
                            port: {
                                name: service.name
                                number: service.port
                            }
                        }
                    ]
                }
            ]
        }
    }
} for service in shared.exposed_services]

ingress_virtual_service_kustomization_resources = [ingress_virtual_service_kustomization_path + ingress_virtual_service.file_name for ingress_virtual_service in ingress_virtual_services]

kustomization = k8s_wrapper.create_kustomization(service_name, [
    ingress_cluster_gateway_file
] + ingress_virtual_service_kustomization_resources)

ingress_cluster_gateway: istio.Gateway = {
    metadata: {
        name: ingress_cluster_gateway_name
    }
    spec: {
        servers: [
            {
                port: {
                    name: "http"
                    number: 80
                    protocol: "HTTP"
                }
                hosts: ["${service.gateway_config.dns_name}.netbird.selfhosted" for service in shared.exposed_services]
            }
        ]
        selector: {
            istio: ingress_cluster_gateway_name
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
## Write ingress-virtual-services
[yaml.dump_to_file(virtual_service.virtual_service, ingress_virtual_service_out_dir) for virtual_service in ingress_virtual_services]
