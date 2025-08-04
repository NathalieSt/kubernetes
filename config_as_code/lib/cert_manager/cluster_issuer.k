import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema Ca:
    secretName: str

schema ClusterIssuerSpec:
    ca?: Ca
    selfSigned?: {str:str}

schema ClusterIssuer:
    apiVersion: "cert-manager.io/v1"
    kind: "ClusterIssuer"
    metadata: meta.ObjectMeta
    spec: ClusterIssuerSpec

