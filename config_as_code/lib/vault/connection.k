import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema VaultConnectionSpec:
    address: str

schema VaultConnection:
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind = "VaultConnection"
    metadata: meta.ObjectMeta
    spec: VaultConnectionSpec

