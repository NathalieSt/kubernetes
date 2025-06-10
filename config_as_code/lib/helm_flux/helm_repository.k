import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema RepositorySpec:
    interval: str
    url: str
    
schema HelmRepository:
    apiVersion: "source.toolkit.fluxcd.io/v1beta2"
    metadata: meta.ObjectMeta
    kind: "HelmRepository"
    spec: RepositorySpec