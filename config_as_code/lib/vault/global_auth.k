import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema VaultAuthGlobalKubernetes:
    audiences: [str]
    mount: str
    role: str
    serviceAccount: str
    tokenExpirationSeconds: int

schema VaultAuthGlobalSpec:
    allowedNamespaces: [str]
    defaultAuthMethod: str
    kubernetes: VaultAuthGlobalKubernetes

schema VaultAuthGlobal:
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind = "VaultAuthGlobal"
    metadata: meta.ObjectMeta
    spec: VaultAuthGlobalSpec

