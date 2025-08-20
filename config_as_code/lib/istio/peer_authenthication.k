import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema Mtls:
    # FIXME: theres a third option, like "none" or something like that
    mode: "STRICT" | "PERMISSIVE"

schema Selector:
    matchLabels: {str:str}

schema PeerAuthenthicationSpec:
    selector?: Selector
    mtls: Mtls

schema PeerAuthenthication:
    apiVersion = "security.istio.io/v1"
    kind = "PeerAuthentication"
    metadata: meta.ObjectMeta
    spec: PeerAuthenthicationSpec

