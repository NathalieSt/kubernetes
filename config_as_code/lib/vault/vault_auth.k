import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema Kubernetes:
    role: str
    serviceAccount: str

schema VaultAuthGlobalRef:
    allowDefault: bool
    namespace: str

schema VaultAuthSpec:
    kubernetes: Kubernetes
    vaultAuthGlobalRef: VaultAuthGlobalRef

schema VaultAuth:
    apiVersion: "secrets.hashicorp.com/v1beta1"
    kind: "VaultAuth"
    metadata: meta.ObjectMeta
    spec: VaultAuthSpec

create_vault_auth = lambda role: str, service_account: str, allow_default: bool, global_auth_namespace: str -> VaultAuth {
    vault_auth = VaultAuth{
        metadata: meta.ObjectMeta{
            name: "vault-auth"
        }
        spec: VaultAuthSpec{
            kubernetes: Kubernetes{
                role: role,
                serviceAccount: service_account
            },
            vaultAuthGlobalRef: VaultAuthGlobalRef{
                allowDefault: allow_default,
                namespace: global_auth_namespace
            }
        }
    }

    vault_auth
}