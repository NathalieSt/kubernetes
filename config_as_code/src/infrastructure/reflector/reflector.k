import k8s.api.core.v1 as core_api
import k8s.api.rbac.v1 as rbac_api
import k8s.apimachinery.pkg.apis.meta.v1 as meta

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ....lib.vault
import ...shared
import yaml

out_dir = shared.cluster_path + "infrastructure/reflector/"
secrets_sub_dir = "secrets/"
secrets_out_dir = out_dir + secrets_sub_dir
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
postgres_creds_file = "postgres-creds.yaml"
vault_auth_file = "vault-auth.yaml"
vault_rbac_file = "vault-rbac.yaml"
# Configs
service_name = "reflector"
chart_url = "https://emberstack.github.io/helm-charts"
chart_name = "reflector"
service_account_name = "reflector-vault-serviceaccount"
# k8s configuration
kustomization: k8s_wrapper.Kustomization = {
    metadata.name = service_name
    resources = [
        namespace_file
        release_file
        repository_file
        # Secrets
        secrets_sub_dir + postgres_creds_file
        secrets_sub_dir + vault_auth_file
        secrets_sub_dir + vault_rbac_file
    ]
}
namespace: core_api.Namespace = k8s_wrapper.create_namespace(service_name, True)

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
                version: "*"
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
        values: {}
    }
}
# Reflector annotations
postgres_creds_reflector_annotations = {
    "reflector.v1.k8s.emberstack.com/reflection-allowed": "true"
    # Control destination namespaces
    "reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "postgres,dawarich,mealie"
    # Auto create reflection for matching namespaces
    "reflector.v1.k8s.emberstack.com/reflection-auto-enabled": "true"
    # Control auto-reflection namespaces
    "reflector.v1.k8s.emberstack.com/reflection-auto-namespaces": "postgres,dawarich,mealie"
}
# Vault configuration
vault_auth: vault.VaultAuth = {
    metadata: {
        name: "vault-auth"
    }
    spec: {
        kubernetes: {
            role: "reflector"
            serviceAccount: service_account_name
        }
        vaultAuthGlobalRef: {
            allowDefault: True
            namespace: shared.global_vault_auth_namespace
        }
    }
}

vault_rbac: vault.VaultRbac = {
    serviceAccount: core_api.ServiceAccount {
        metadata: meta.ObjectMeta {
            name: service_account_name
        }
    }
    role: rbac_api.Role {
        metadata: meta.ObjectMeta {
            name: "vault-token-reviewer"
        }
        rules: [
            rbac_api.PolicyRule {
                apiGroups: ["authentication.k8s.io"]
                resources: ["tokenreviews"]
                verbs: ["create"]
            }
            rbac_api.PolicyRule {
                apiGroups: ["authorization.k8s.io"]
                resources: ["subjectaccessreviews"]
                verbs: ["create"]
            }
        ]
    }
    rolebinding: rbac_api.RoleBinding {
        metadata: meta.ObjectMeta {
            name: "reflector-vault-binding"
        }
        subjects: [
            rbac_api.Subject {
                kind: "ServiceAccount"
                name: service_account_name
            }
        ]
        roleRef: rbac_api.RoleRef {
            kind: "Role"
            name: "vault-token-reviewer"
            apiGroup: "rbac.authorization.k8s.io"
        }
    }
}

postgres_creds: vault.VaultStaticSecret = {
    metadata: {
        name: "postgres-creds"
    }
    spec: {
        vaultAuthRef: vault_auth.metadata.name
        mount: "kvv2"
        type: "kv-v2"
        path: "postgres"
        refreshAfter: "10s"
        destination: {
            create: True
            name: shared.postgres_creds_secret_name
            annotations: postgres_creds_reflector_annotations
        }
    }
}
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
# Write secrets to files
yaml.dump_to_file(postgres_creds, secrets_out_dir + postgres_creds_file)
yaml.dump_to_file(vault_auth, secrets_out_dir + vault_auth_file)
yaml.dump_all_to_file([vault_rbac[key] for key in vault_rbac], secrets_out_dir + vault_rbac_file)
