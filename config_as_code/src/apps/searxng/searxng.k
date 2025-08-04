import ....lib.k8s_wrapper
import ....lib.vault
import ...shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "apps/searxng/"
secret_dir = "secrets/"
secret_out_dir = out_dir + secret_dir
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
configmap_file = "configmap.yaml"
deployment_file = "deployment.yaml"
service_file = "service.yaml"
vault_rbac_file = "vault_rbac.yaml"
vault_auth_file = "vault_auth.yaml"
vault_static_secret_file = "vault_static_secret.yaml"
# Configs
service_name = "searxng"
label_app_name = "searxng"
k8s_service_name = "searxng"
configmap_name = "searxng-configmap"
deployment_configmap_volume = "searxng-config"
secret_name = "searxng-vpn"
vault_role = "vpn"
secret_path = "vpn"
port = 8080
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    configmap_file
    deployment_file
    service_file
    # secrets
    secret_dir + vault_rbac_file
    secret_dir + vault_auth_file
    secret_dir + vault_static_secret_file
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
    # secret_key is overwritten by ${SEARXNG_SECRET}, currently a random string, later an env variable will replace it
    # URL to connect valkey database is overwritten by ${SEARXNG_VALKEY_URL}
    # FIXME: maybe extract the inline yaml into a separate file and use a kcl function to turn it into a string here??? Could be useful for nginx gateway later
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
                secret_key: "mBrieVOfZQzc7"
                limiter: true
                image_proxy: true

            valkey:
                url: valkey://valkey.valkey.svc.cluster.local:6379/0
        """
    }
}

vault_rbac = vault.create_vault_rbac(service_name)

vault_auth: vault.VaultAuth = {
    metadata: {
        name: "vault-auth"
    }
    spec: {
        kubernetes: {
            role: vault_role
            serviceAccount: vault_rbac.serviceAccount.metadata.name
        }
        vaultAuthGlobalRef: {
            allowDefault: True
            namespace: shared.global_vault_auth_namespace
        }
    }
}

vault_static_secret: vault.VaultStaticSecret = {
    metadata: {
        name: secret_name
    }
    spec: {
        vaultAuthRef: vault_auth.metadata.name
        mount: "kvv2"
        type: "kv-v2"
        path: secret_path
        refreshAfter: "10s"
        destination: {
            create: True
            name: secret_name
        }
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
                    {
                        name: "gluetun"
                        image: "qmcgaw/gluetun:v3.40"
                        securityContext: {
                            capabilities: {
                                add: [
                                    "NET_ADMIN"
                                ]
                            }
                        }
                        env: [
                            {
                                name: "VPN_SERVICE_PROVIDER"
                                value: "protonvpn"
                            }
                            {
                                name: "VPN_TYPE"
                                value: "wireguard"
                            }
                            {
                                name: "WIREGUARD_PRIVATE_KEY"
                                valueFrom: {
                                    secretKeyRef: {
                                        name: secret_name
                                        key: "key"
                                    }
                                }
                            }
                            {
                                name: "DOT"
                                value: "on"
                            }
                            {
                                name: "DNS_KEEP_NAMESERVER"
                                value: "on"
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
# secrets
yaml.dump_all_to_file([vault_rbac[key] for key in vault_rbac], secret_out_dir + vault_rbac_file)
yaml.dump_to_file(vault_auth, secret_out_dir + vault_auth_file)
yaml.dump_to_file(vault_static_secret, secret_out_dir + vault_static_secret_file)
