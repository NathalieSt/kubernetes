import ....lib.k8s_wrapper
import ....lib.vault
import ....lib.cert_manager
import ....lib.nginx
import k8s.api.core.v1 as core_api
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import ...shared
import yaml

out_dir = shared.cluster_path + "infrastructure/nginx-gateway/"
secrets_subdir = "secrets/"
secrets_out_dir = out_dir + secrets_subdir
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
deployment_file = "deployment.yaml"
certificate_file = "nginx-gateway-certificate.yaml"
configmap_file = "nginx-gateway-config.yaml"
k8s_service_file = "service.yaml"
## Vault specific files
vault_rbac_file = "vault-rbac.yaml"
vault_auth_file = "vault-auth.yaml"
vault_netbird_static_secret_file = "vault-netbird-static-secret.yaml"
vault_nginx_static_secret_file = "vault-nginx-static-secret.yaml"
# Configs
service_name = shared.nginx_gateway.name
label_app_name = service_name
port = 80
k8s_service_name = service_name
configmap_name = "nginx-gateway-config"
certificate_name = "nginx-gateway-tls"
## Vault specific configs
vault_role = "nginx"
netbird_secret_name = "netbird-setup-key-vault-secret"
netbird_secret_path = "netbird/setup-key"
# FIXME: change to something like "nginx user secret"
nginx_secret_name = "nginx-secret"
nginx_secret_path = "nginx"
nginx_secret_volume = "nginx-secret-volume"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    deployment_file
    certificate_file
    configmap_file
    k8s_service_file
    # Vault
    secrets_subdir + vault_rbac_file
    secrets_subdir + vault_auth_file
    secrets_subdir + vault_netbird_static_secret_file
    secrets_subdir + vault_nginx_static_secret_file
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

vault_netbird_static_secret: vault.VaultStaticSecret = {
    metadata: {
        name: netbird_secret_name
    }
    spec: {
        vaultAuthRef: vault_auth.metadata.name
        mount: "kvv2"
        type: "kv-v2"
        path: netbird_secret_path
        refreshAfter: "10s"
        destination: {
            create: True
            name: netbird_secret_name
        }
    }
}

vault_nginx_static_secret: vault.VaultStaticSecret = {
    metadata: {
        name: nginx_secret_name
    }
    spec: {
        vaultAuthRef: vault_auth.metadata.name
        mount: "kvv2"
        type: "kv-v2"
        path: nginx_secret_path
        refreshAfter: "10s"
        destination: {
            create: True
            name: nginx_secret_name
            transformation: {
                templates: {
                    "htpasswd": {
                        text: """
                        {{- $user := get .Secrets "user" -}}
                        {{- $password := get .Secrets "password" -}}
                        {{- printf "%s:%s" $user $password -}}              
                        """
                    }
                }
            }
        }
    }
}

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
        "default.conf": nginx.create_conf_from_services(shared.exposed_services)
    }
}

certificate: cert_manager.Certificate = {
    apiVersion: "cert-manager.io/v1"
    kind: "Certificate"
    metadata: {
        name: certificate_name
    }
    spec: {
        subject: {
            organizations: [
                "cert-manager"
            ]
        }
        secretName: certificate_name
        issuerRef: {
            name: "my-ca-issuer"
            kind: "ClusterIssuer"
        }
        dnsNames: [
            "*.cluster.netbird.selfhosted"
        ]
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: service_name
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
                        name = "nginx"
                        image = "nginx:alpine"
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
                        ports = [
                            {
                                containerPort = 80
                            }
                            {
                                containerPort = 443
                            }
                        ]
                        volumeMounts = [
                            {
                                name = "nginx-config"
                                mountPath = "/etc/nginx/conf.d"
                                subPath: "default.conf"
                            }
                            {
                                name = "nginx-tls"
                                mountPath = "/etc/nginx/ssl"
                                readOnly = True
                            }
                            {
                                name = "shared-tmp"
                                mountPath = "/tmp"
                            }
                            {
                                name: nginx_secret_volume
                                mountPath: "/etc/nginx/conf.d/nginx.htpasswd"
                                subPath: "htpasswd"
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
                                # Old hardcoded string: "jellyfin.cluster,vault.cluster,gatus.cluster,wikijs.cluster,searxng.cluster"
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
                volumes = [
                    {
                        name = "nginx-config"
                        configMap = {
                            name = configmap_name
                        }
                    }
                    {
                        name = "nginx-tls"
                        secret = {
                            secretName = certificate_name
                            defaultMode = 256
                        }
                    }
                    {
                        name = "shared-tmp"
                        emptyDir = {}
                    }
                    {
                        name: nginx_secret_volume
                        secret: {
                            secretName: nginx_secret_name
                        }
                    }
                ]
            }
        }
    }
}

k8s_service: core_api.Service = {
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
yaml.dump_to_file(certificate, out_dir + certificate_file)
yaml.dump_to_file(k8s_service, out_dir + k8s_service_file)
## Write vault configuration to files
yaml.dump_all_to_file([vault_rbac[key] for key in vault_rbac], secrets_out_dir + vault_rbac_file)
yaml.dump_to_file(vault_auth, secrets_out_dir + vault_auth_file)
yaml.dump_to_file(vault_netbird_static_secret, secrets_out_dir + vault_netbird_static_secret_file)
yaml.dump_to_file(vault_nginx_static_secret, secrets_out_dir + vault_nginx_static_secret_file)
