import k8s.apimachinery.pkg.apis.meta.v1 as meta
import ..shared

schema Server:
    port: shared.IstioPort
    hosts: [str]

schema GatewaySelector:
    istio: str

schema GatewaySpec:
    servers: [Server]
    selector: GatewaySelector

schema Gateway:
    apiVersion = "networking.istio.io/v1"
    kind = "Gateway"
    metadata: meta.ObjectMeta
    servers: [Server]

