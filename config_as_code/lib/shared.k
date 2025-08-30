import .flux

schema VirtualServiceConfig:
    client_protocol?: "http" | "https"

schema GatewayConfig:
    dns_name: str
    requires_websockets: bool
    requires_larger_body: bool
    extra_server_config?: str
    extra_location_config?: str
    are_headers_forwarded?: bool

schema ClusterServiceFluxKustomization:
    kustomization_path: str
    dependencies: [flux.FluxKustomizationDependency]

schema KedaScaling:
    start: str
    end: str
    desiredReplicas: int

schema ClusterService:
    name: str
    namespace: str
    type: "app" | "infrastructure" | "istio" | "monitoring"
    cluster_url?: str
    port?: int
    gateway_config?: GatewayConfig
    virtual_service_config?: VirtualServiceConfig
    flux_kustomization: ClusterServiceFluxKustomization
    keda_scaling?: KedaScaling

schema ExternalService:
    name: str
    url: str

schema IstioPort:
    number: int
    name: str
    protocol: "TCP" | "HTTP"

