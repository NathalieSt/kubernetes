import ....lib.k8s_wrapper
import ...shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml
import outdent

out_dir = shared.cluster_path + "apps/searxng/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
configmap_file = "configmap.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
service_entry_file = "service-entry.yaml"
# Configs
service_name = "searxng"
label_app_name = "searxng"
k8s_service_name = "searxng"
configmap_name = "searxng-configmap"
image_version = "2025.8.3-2e62eb5"
deployment_configmap_volume = "searxng-config"
port = 8080
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    configmap_file
    deployment_file
    k8s_service_file
    service_entry_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)

configmap = core_api.ConfigMap {
    apiVersion: "v1"
    kind: "ConfigMap"
    metadata: {
        name: configmap_name
        annotations: {
            "app.kubernetes.io/name": label_app_name
        }
    }
    # secret_key is overwritten by ${SEARXNG_SECRET}, currently a random string, later an env variable will replace it
    # URL to connect valkey database is overwritten by ${SEARXNG_VALKEY_URL}
    # FIXME: maybe extract the inline yaml into a separate file and use a kcl function to turn it into a string here??? Could be useful for nginx gateway later
    data: {
        "settings.yml": outdent.outdent("""
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
                secret_key: "mBrieVOfZQzc7"
                limiter: true
                image_proxy: true

            valkey:
                url: valkey://valkey.valkey.svc.cluster.local:6379/0

            outgoing:
                request_timeout: 2.0       # default timeout in seconds, can be override by engine
                max_request_timeout: 10.0  # the maximum timeout in seconds
                useragent_suffix: ""       # information like an email address to the administrator
                pool_connections: 100      # Maximum number of allowable connections, or null
                                            # for no limits. The default is 100.
                pool_maxsize: 10           # Number of allowable keep-alive connections, or null
                                            # to always allow. The default is 10.
                enable_http2: true         # See https://www.python-httpx.org/http2/
                proxies:
                    all://:
                        - http://gluetun-proxy.gluetun-proxy.svc.cluster.local:8888
        """)
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: "searxng-deployment"
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
                        name: "searxng"
                        image: "searxng/searxng:${image_version}"
                        env: [
                            {
                                name: "NO_PROXY"
                                value: "localhost,127.0.0.1,.svc,.cluster.local,gluetun-proxy.gluetun-proxy.svc.cluster.local"
                            }
                        ]
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
            "app.kubernetes.io/version": image_version
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
yaml.dump_to_file(service, out_dir + k8s_service_file)
