forgejo_kustomization = {
    apiVersion = "kustomize.config.k8s.io/v1beta1"
    kind = "Kustomization"
    metadata = {
        name = "forgejo"
    }
    resources = [
        "repository.yaml"
        "release.yaml"
        "ingress.yaml"
        "namespace.yaml"
    ]
}
