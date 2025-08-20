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

