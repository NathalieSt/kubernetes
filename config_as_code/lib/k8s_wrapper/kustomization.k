import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema Kustomization:
    apiVersion: "kustomize.config.k8s.io/v1beta1"
    kind: "Kustomization"
    metadata: meta.ObjectMeta
    resources: [str]

create_kustomization = lambda name: str, resources: [str] -> Kustomization {
    {
        metadata.name = name
        resources = resources
    }
}