import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema FluxKustomizationSourceRef:
    # TODO: what are the other options?
    kind: "GitRepository"
    name: str

schema FluxKustomizationDependency:
    name: str

schema FluxKustomizationSpec:
    interval: str
    targetNamespace: str
    sourceRef: FluxKustomizationSourceRef
    path: str
    prune: bool
    wait: bool
    timeout: str
    dependsOn: [FluxKustomizationDependency]

schema FluxKustomization:
    apiVersion: "kustomize.toolkit.fluxcd.io/v1"
    kind: "Kustomization"
    metadata: meta.ObjectMeta
    spec: FluxKustomizationSpec

