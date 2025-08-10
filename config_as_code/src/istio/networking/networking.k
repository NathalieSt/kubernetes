import ....lib.istio
import ....lib.k8s_wrapper
import ...shared
import yaml

schema VirtualServiceAndFile:
    file_name: str
    virtual_service: istio.VirtualService

# Output file paths
out_dir = shared.cluster_path + "istio/networking/"
ingress_out_dir = out_dir + "ingress/"
ingress_virtual_service_out_dir = ingress_out_dir + "virtual-services/"
# Kustomization file paths
ingress_kustomization_path = "ingress/"
ingress_virtual_service_kustomization_path = ingress_kustomization_path + "virtual-services/"
# Filenames
kustomization_file = "kustomization.yaml"
ingress_cluster_gateway_file = "ingress-cluster-gateway.yaml"
# Configs
service_name = "istio-networking"
ingress_cluster_gateway_name = "ingress-cluster-gateway"
ingress_name = "ingress"
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
                    route: [
                        {
                            destination: {
                                host: "${service.cluster_url}"
                                port: {
                                    # FIXME: Include port name in shared.<services>.gateway_config
                                    name: service.name
                                    number: service.port
                                    # FIXME: Include protocol in shared.<services>.gateway_config
                                    protocol: "HTTP"
                                }
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
    ingress_kustomization_path + ingress_cluster_gateway_file
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
            istio: ingress_name
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(ingress_cluster_gateway, ingress_out_dir + ingress_cluster_gateway_file)
## Write ingress-virtual-services
_ = [yaml.dump_to_file(virtual_service.virtual_service, ingress_virtual_service_out_dir + virtual_service.file_name) for virtual_service in ingress_virtual_services]
