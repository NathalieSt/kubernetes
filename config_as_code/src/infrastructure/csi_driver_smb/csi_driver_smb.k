import ....lib.k8s_wrapper
import ....lib.helm_flux
import ....lib.vault
import ...shared
import yaml

out_dir = shared.cluster_path + "infrastructure/csi-driver-smb/"
secrets_path = out_dir + "secrets/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
storage_class_file = "storage_class.yaml"
vault_rbac_file = "vault_rbac.yaml"
vault_auth_file = "vault_auth.yaml"
vault_static_secret_file = "vault_static_secret.yaml"
# Configs
service_name = "csi-driver-smb"
chart_url = "https://raw.githubusercontent.com/kubernetes-csi/csi-driver-smb/master/charts"
chart_name = "csi-driver-smb"
vault_role = "smb"
secret_name = "smbcreds"

secrets_subdir = "secrets/"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
    storage_class_file
    secrets_subdir + vault_rbac_file
    secrets_subdir + vault_auth_file
    secrets_subdir + vault_static_secret_file
])
namespace = k8s_wrapper.create_namespace(service_name)
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

storage_class = {
    apiVersion: "storage.k8s.io/v1"
    kind: "StorageClass"
    metadata.name = shared.smb_storage_class_name
    provisioner: "smb.csi.k8s.io"
    parameters: {
        source: "//u465810.your-storagebox.de/backup"
        "csi.storage.k8s.io/provisioner-secret-name": secret_name
        "csi.storage.k8s.io/provisioner-secret-namespace": service_name
        "csi.storage.k8s.io/node-stage-secret-name": secret_name
        "csi.storage.k8s.io/node-stage-secret-namespace": service_name
    }
    reclaimPolicy: "Delete"
    volumeBindingMode: "Immediate"
    allowVolumeExpansion: True
    mountOptions: ["dir_mode=0755", "file_mode=0644"]
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
        path: "smbcreds"
        refreshAfter: "10s"
        destination: {
            create: True
            name: secret_name
        }
    }
}
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
# Write storage class to file
yaml.dump_to_file(storage_class, out_dir + storage_class_file)
# Write vault configuration to files
yaml.dump_all_to_file([vault_rbac[key] for key in vault_rbac], secrets_path + vault_rbac_file)
yaml.dump_to_file(vault_auth, secrets_path + vault_auth_file)
yaml.dump_to_file(vault_static_secret, secrets_path + vault_static_secret_file)
