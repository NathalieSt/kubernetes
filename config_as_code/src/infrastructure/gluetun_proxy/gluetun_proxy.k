import ....lib.k8s_wrapper
import ...shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml
import ....lib.vault

out_dir = shared.cluster_path + "infrastructure/gluetun-proxy/"
# Vault secrets specific path
secret_dir = "secrets/"
secret_out_dir = out_dir + secret_dir
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
deployment_file = "deployment.yaml"
service_file = "service.yaml"
## Vault specific filenames
vault_rbac_file = "vault_rbac.yaml"
vault_auth_file = "vault_auth.yaml"
vault_static_secret_file = "vault_static_secret.yaml"
# Configs
service_name = "gluetun-proxy"
label_app_name = "gluetun-proxy"
k8s_service_name = "gluetun-proxy"
port = 8888
## Vault specific configs
secret_name = "gluetun-proxy-vpn"
vault_role = "vpn"
secret_path = "vpn"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    deployment_file
    service_file
    # Vault specific files
    secret_dir + vault_rbac_file
    secret_dir + vault_auth_file
    secret_dir + vault_static_secret_file
])
namespace = k8s_wrapper.create_namespace(service_name)

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
        name: "gluetun-proxy-deployment"
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
                        name: "gluetun"
                        image: "qmcgaw/gluetun:v3.40"
                        ports: [
                            {
                                name: "gluetun-proxy"
                                containerPort: port
                            }
                        ]
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
                                value: "openvpn"
                            }
                            {
                                name: "OPENVPN_USER"
                                valueFrom: {
                                    secretKeyRef: {
                                        name: secret_name
                                        key: "user"
                                    }
                                }
                            }
                            {
                                name: "OPENVPN_PASSWORD"
                                valueFrom: {
                                    secretKeyRef: {
                                        name: secret_name
                                        key: "password"
                                    }
                                }
                            }
                            {
                                name: "SERVER_COUNTRIES"
                                value: "Austria,Finland,Malta"
                            }
                            {
                                name: "FIREWALL_INPUT_PORTS"
                                value: "8888"
                            }
                            {
                                name: "FIREWALL_OUTBOUND_SUBNETS"
                                # Allow local cluster traffic
                                value: "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,10.42.0.0/15"
                            }
                            {
                                name: "DNS_KEEP_NAMESERVER"
                                value: "on"
                            }
                            {
                                name: "HTTPPROXY"
                                value: "on"
                            }
                            {
                                name: "HTTPPROXY_STEALTH"
                                value: "on"
                            }
                        ]
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
                name: "gluetun-proxy"
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
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(service, out_dir + service_file)
## secrets
yaml.dump_all_to_file([vault_rbac[key] for key in vault_rbac], secret_out_dir + vault_rbac_file)
yaml.dump_to_file(vault_auth, secret_out_dir + vault_auth_file)
yaml.dump_to_file(vault_static_secret, secret_out_dir + vault_static_secret_file)
