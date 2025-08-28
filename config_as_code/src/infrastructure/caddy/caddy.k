import ....lib.k8s_wrapper
import ....lib.shared as lib_shared
import ...shared
import ...versions
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml
import outdent

out_dir = shared.cluster_path + "infrastructure/caddy/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
configmap_file = "configmap.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
cert_pvc_file = "caddy-cert-pvc.yaml"
# Configs
service_name = "caddy"
label_app_name = "caddy"
k8s_service_name = "caddy"
configmap_name = "caddy-configmap"
deployment_configmap_volume = "caddy-config"
image_version = versions.docker_images.caddy
netbird_secret_name = shared.netbird_secret_name
port = 80
cert_pvc_name = "caddy-cert-pvc"
cert_pvc_volume = "caddy-cert-pvc-volume"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    configmap_file
    deployment_file
    k8s_service_file
    cert_pvc_file
])
namespace = k8s_wrapper.create_namespace(service_name, False)

cert_pvc: core_api.PersistentVolumeClaim = {
    apiVersion: "v1"
    kind: "PersistentVolumeClaim"
    metadata: {
        name: cert_pvc_name
    }
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_local_storage_class_v2_name
        resources.requests.storage: "1Gi"
    }
}

# FIXME: put into lib.caddy
get_websocket_support_if_required = lambda service: lib_shared.ClusterService -> str {
    result = ""
    if service.gateway_config.requires_websockets == True:
        result = """
        @websockets {
            header Connection *Upgrade*
            header Upgrade    websocket
        }
        reverse_proxy @websockets ${shared.istio_gateway_ip}:80 {
            header_up Host {host}
        }
    """

    result
}

forward_headers_if_required = lambda service: lib_shared.ClusterService -> str {
    result = ""
    if service.gateway_config.are_headers_forwarded == True:
        result = """
            header_up X-Forwarded-Proto {scheme}
            header_up X-Forwarded-Host {host}
            header_up X-Forwarded-Uri {uri}
    """

    result
}

caddy_services = ("").join(["""
    ${service.gateway_config.dns_name}.netbird.selfhosted:443 {
        tls internal
        reverse_proxy ${shared.istio_gateway_ip}:80 {
            header_up Host {host}
            ${forward_headers_if_required(service)}
        }

        ${get_websocket_support_if_required(service)}
    }
""" for service in shared.exposed_services])

configmap = core_api.ConfigMap {
    apiVersion: "v1"
    kind: "ConfigMap"
    metadata: {
        name: configmap_name
        annotations: {
            "app.kubernetes.io/name": label_app_name
        }
    }
    data: {
        "Caddyfile": outdent.outdent("""
            ${caddy_services}
        """)
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: "caddy-deployment"
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        replicas: 1
        selector: {
            matchLabels: {
                "app.kubernetes.io/name": label_app_name
                "app.kubernetes.io/version": image_version
            }
        }
        template: {
            metadata: meta.ObjectMeta {
                labels: {
                    "app.kubernetes.io/name": label_app_name
                    "app.kubernetes.io/version": image_version
                }
            }
            spec = core_api.PodSpec {
                containers: [
                    {
                        name: "caddy"
                        image: "caddy:${image_version}"
                        ports: [
                            {
                                name: "https"
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                mountPath: "/etc/caddy/"
                                name: deployment_configmap_volume
                            }
                            {
                                mountPath: "/data/caddy/pki/authorities/local/"
                                name: cert_pvc_volume
                            }
                        ]
                    }
                    {
                        name = "netbird-agent"
                        image = "netbirdio/netbird:latest"
                        env = [
                            {
                                name = "NB_SETUP_KEY"
                                valueFrom = {
                                    secretKeyRef = {
                                        name = netbird_secret_name
                                        key = "setup-key"
                                    }
                                }
                            }
                            {
                                name = "NB_HOSTNAME"
                                value = "cluster"
                            }
                            {
                                name = "NB_MANAGEMENT_URL"
                                value = "https://netbird.nathalie-stiefsohn.eu"
                            }
                            {
                                name = "NB_EXTRA_DNS_LABELS"
                                # Old hardcoded string: "gatus.cluster,searxng.cluster,kiali.cluster"
                                value = ",".join([service.gateway_config.dns_name for service in shared.exposed_services])
                            }
                        ]
                        resources = {
                            requests = {
                                cpu = "50m"
                                memory = "64Mi"
                            }
                            limits = {
                                cpu = "100m"
                                memory = "128Mi"
                            }
                        }
                        securityContext = {
                            privileged = True
                        }
                    }
                ]
                volumes: [
                    {
                        name: deployment_configmap_volume
                        configMap: {
                            name: configmap_name
                        }
                    }
                    {
                        name: cert_pvc_volume
                        persistentVolumeClaim: {
                            claimName: cert_pvc_name
                        }
                    }
                ]
            }
        }
    }
}

service: core_api.Service = {
    apiVersion: "v1"
    kind: "Service"
    metadata: {
        name: k8s_service_name
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        ports: [
            {
                name: "http"
                port: port
                targetPort: port
            }
        ]
        selector: {
            "app.kubernetes.io/name": label_app_name
        }
    }
}
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(configmap, out_dir + configmap_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(service, out_dir + k8s_service_file)
yaml.dump_to_file(cert_pvc, out_dir + cert_pvc_file)
