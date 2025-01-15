import .variables

homepage_kustomization = {
    apiVersion = "kustomize.config.k8s.io/v1beta1"
    kind = "Kustomization"
    metadata = {
        name = "homepage"
    }
    resources = [
        "rbac.yaml"
        "deployment.yaml"
        "service.yaml"
        "configmap.yaml"
        "ingress.yaml"
    ]
}
