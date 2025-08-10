import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema VirtualServicePort:
    number: int

schema Destination:
    host: str
    port: VirtualServicePort

schema Set:
    "X-Forwarded-Proto": "http" | "https"

schema Request:
    set: Set

schema Headers:
    request: Request

schema Route:
    destination: Destination
    headers?: Headers

schema Http:
    route: [Route]

schema VirtualServiceSpec:
    hosts: [str]
    gateways: [str]
    http: [Http]

schema VirtualService:
    apiVersion = "networking.istio.io/v1"
    kind = "VirtualService"
    metadata: meta.ObjectMeta
    spec: VirtualServiceSpec

