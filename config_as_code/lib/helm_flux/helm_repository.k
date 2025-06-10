import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema RepositorySpec:
    interval: "24h"
    url: str
    
schema HelmRepository:
    apiVersion: "source.toolkit.fluxcd.io/v1beta2"
    metadata: meta.ObjectMeta
    kind: "HelmRepository"
    spec: RepositorySpec

create_helm_repository = lambda url: str, repo_name: str -> HelmRepository {
    repo = {   
        metadata.name = repo_name
        spec.url= url
    }

    repo
}