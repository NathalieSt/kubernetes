import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema VirtualServicePort:
    number: int

schema Destination:
    host: str
    port: VirtualServicePort

schema Route:
    destination: Destination

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

