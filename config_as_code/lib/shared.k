schema GatusConfig:
    endpointType: "http" | "tcp"
    conditions?: [str]

schema GatewayConfig:
    dns_name: str
    requires_websockets: bool
    requires_larger_body: bool
    extra_server_config?: str
    extra_location_config?: str

schema ClusterService:
    name: str
    namespace: str
    cluster_url: str
    port: int
    gatus_config?: GatusConfig
    gateway_config?: GatewayConfig

schema ExternalService:
    name: str
    url: str

