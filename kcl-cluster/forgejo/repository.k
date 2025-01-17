forgejo_repository = {
    apiVersion: "source.toolkit.fluxcd.io/v1beta2"
    kind: "OCIRepository"
    metadata: {
        name: "forgejo"
    }
    spec: {
        interval: "50m0s"
        url: "oci://code.forgejo.org/forgejo-helm/forgejo"
        ref: {
            semver: "11.0.1"
        }
    }
}
