import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema ClusterStorage:
    storageClass: str
    size: str

schema ClusterSpec:
    instances: int
    storage: ClusterStorage

schema Cluster:
    apiVersion = "postgresql.cnpg.io/v1"
    kind = "Cluster"
    metadata: meta.ObjectMeta
    spec: ClusterSpec

