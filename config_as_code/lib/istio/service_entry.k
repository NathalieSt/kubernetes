import k8s.apimachinery.pkg.apis.meta.v1 as meta
import ..shared

schema ServiceEntrySpec:
    hosts: [str]
    ports: [shared.IstioPort]
    location: "MESH_EXTERNAL" | "MESH_INTERNAL"
    resolution: "DNS" | "STATIC"

schema ServiceEntry:
    apiVersion = "networking.istio.io/v1alpha3"
    kind = "ServiceEntry"
    metadata: meta.ObjectMeta
    spec: ServiceEntrySpec

