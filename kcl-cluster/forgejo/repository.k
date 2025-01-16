forgejo_repository = {
    apiVersion: "source.toolkit.fluxcd.io/v1"
    kind: "HelmRepository"
    metadata: {
        name: "forgejo"
    }
    spec: {
        type: "oci"
        interval: "50m0s"
        url: "oci://code.forgejo.org/forgejo-helm/forgejo"
        ref: {
            tag: "11.0.1"
        }
    }
}
