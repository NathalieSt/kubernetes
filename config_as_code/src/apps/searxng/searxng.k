import ....lib.k8s_wrapper
import ...shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "apps/searxng/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
configmap_file = "configmap.yaml"
deployment_file = "deployment.yaml"
service_file = "service.yaml"
# Configs
service_name = "searxng"
label_app_name = "searxng"
k8s_service_name = "searxng"
configmap_name = "searxng-configmap"
deployment_configmap_volume = "searxng-config"
port = 8888
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    configmap_file
    deployment_file
    service_file
])
namespace = k8s_wrapper.create_namespace(service_name)

configmap = core_api.ConfigMap {
    apiVersion: "v1"
    kind: "ConfigMap"
    metadata: {
        name: configmap_name
        annotations: {
            "app.kubernetes.io/name": label_app_name
        }
    }
    # secret_key is overwritten by ${SEARXNG_SECRET}
    # URL to connect valkey database is overwritten by ${SEARXNG_VALKEY_URL}
    data: {
        "settings.yml": """
            # SearXNG settings
            use_default_settings: true

            general:
                debug: false
                instance_name: "Nathalies SearXNG"

            search:
                safe_search: 2
                autocomplete: 'duckduckgo'
                formats:
                    - html

            server:
                secret_key: "ultrasecretkey"
                limiter: true
                image_proxy: true

            valkey:
                url: valkey://valkey.valkey.svc.cluster.local:6379/0
        """
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: "searxng-deployment"
        labels: {
            "app.kubernetes.io/name": label_app_name
        }
    }
    spec: {
        replicas: 1
        selector: {
            matchLabels: {
                "app.kubernetes.io/name": label_app_name
            }
        }
        template: {
            metadata: meta.ObjectMeta {
                labels: {
                    "app.kubernetes.io/name": label_app_name
                }
            }
            spec = core_api.PodSpec {
                containers: [
                    {
                        name: "searxng"
                        image: "searxng/searxng:2025.8.3-2e62eb5"
                        ports: [
                            {
                                name: "searxng"
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                mountPath: "/etc/searxng"
                                name: deployment_configmap_volume
                            }
                        ]
                    }
                ]
                volumes: [
                    {
                        name: deployment_configmap_volume
                        configMap: {
                            name: configmap_name
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
        }
    }
    spec: {
        ports: [
            {
                name: "searxng"
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
yaml.dump_to_file(service, out_dir + service_file)
