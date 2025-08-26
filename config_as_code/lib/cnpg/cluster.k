import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema ClusterStorage:
    storageClass: str
    size: str

schema InheritedMetadata:
    annotations: {str:str}

schema ClusterSpec:
    instances: int
    storage: ClusterStorage
    inheritedMetadata?: InheritedMetadata

schema Cluster:
    apiVersion = "postgresql.cnpg.io/v1"
    kind = "Cluster"
    metadata: meta.ObjectMeta
    spec: ClusterSpec

