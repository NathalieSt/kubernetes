import yaml
import ....lib.k8s_wrapper
import ....lib.helm_flux
import ....lib.vault
import ...shared
import ...versions

out_dir = shared.cluster_path + "infrastructure/vault-secrets-operator/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
vault_rbac_file = "vault-rbac.yaml"
vault_auth_global_file = "vault-auth-global.yaml"
vault_connection_file = "vault-connection.yaml"
# Configs
service_name = "vault-secrets-operator"
chart_url = "https://helm.releases.hashicorp.com"
chart_name = "vault-secrets-operator"
vault_auth_global_name = "global-vault-auth"

# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
    vault_rbac_file
    vault_auth_global_file
    vault_connection_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)

vault_rbac = vault.create_vault_rbac(service_name)

vault_auth_global: vault.VaultAuthGlobal = {
    metadata: {
        name: vault_auth_global_name
    }
    spec: {
        allowedNamespaces: ["reflector"]
        defaultAuthMethod: "kubernetes"
        kubernetes: {
            audiences: ["vault"]
            mount: "kubernetes"
            role: vault_rbac.role.metadata.name
            serviceAccount: vault_rbac.serviceAccount.metadata.name
            tokenExpirationSeconds: 600
        }
    }
}

vault_connection: vault.VaultConnection = {
    metadata: {
        name: "default"
    }
    spec: {
        address: "http://${shared.vault.cluster_url}:8200"
    }
}

# Helm configuration
repo: helm_flux.HelmRepository = {
    metadata.name = service_name
    kind: "HelmRepository"
    spec: {
        interval: "24h"
        url: chart_url
    }
}
release: helm_flux.HelmRelease = {
    metadata.name = service_name
    spec: {
        releaseName: service_name
        chart: {
            spec: {
                chart: chart_name
                version: versions.helm_charts.vault_secrets_operator
                interval: "24h"
                sourceRef: {
                    kind: "HelmRepository"
                    name: service_name
                }
            }
        }
        install.remediation.retries: 3
        timeout: "5m"
        interval: "24h"
        values: {
            controller: {
                annotations: {
                    "traffic.sidecar.istio.io/excludeOutboundPorts": "8200"
                }
            }
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
## Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
## Write vault connection config to files
yaml.dump_all_to_file([vault_rbac[key] for key in vault_rbac], out_dir + vault_rbac_file)
yaml.dump_to_file(vault_connection, out_dir + vault_connection_file)
yaml.dump_to_file(vault_auth_global, out_dir + vault_auth_global_file)
