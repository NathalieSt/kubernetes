import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema DatabaseCluster:
    name: str

schema DatabaseSpec:
    cluster: DatabaseCluster
    owner: str

schema Database:
    apiVersion = "postgresql.cnpg.io/v1"
    kind = "Database"
    metadata: meta.ObjectMeta
    spec: DatabaseSpec

