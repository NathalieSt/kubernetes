import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema TLS:
    mode: "DISABLE" | "SIMPLE" | "MUTUAL" | "ISTIO_MUTUAL"

schema TrafficPolicy:
    tls: TLS

schema DestinationRuleSpec:
    host: str
    trafficPolicy: TrafficPolicy

schema DestinationRule:
    apiVersion = "networking.istio.io/v1"
    kind = "DestinationRule"
    metadata: meta.ObjectMeta
    spec: DestinationRuleSpec

