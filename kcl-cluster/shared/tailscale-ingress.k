schema TailscaleIngress:
    apiVersion: "networking.k8s.io/v1" = "networking.k8s.io/v1"
    kind: "Ingress" = "Ingress"
    metadata: Metadata
    spec: Spec

schema Metadata:
    name: str

schema Spec:
    ingressClassName: "tailscale" = "tailscale"
    rules: [Http]
    tls: [Tls]

schema Http:
    http: Paths

schema Paths:
    paths: [Path]

schema Path:
    path: str
    pathType: "Prefix" = "Prefix"
    backend: Backend

schema Backend:
    service: Service

schema Service:
    name: str
    port: Port

schema Port:
    number: int

schema Tls:
    hosts: [str]

