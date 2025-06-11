schema Netbird:
    setup_key: str
    key: str

schema GatusConfig:
    endpointType: "http" | "tcp"
    conditions?: [str]

schema ClusterService:
    name: str
    namespace: str
    cluster_url: str
    gatus: GatusConfig

cluster_path = "../../../../cluster/"
global_vault_auth_namespace = "vault-secrets-operator"
ca_issuer = "my-ca-issuer"
netbird: Netbird = {
    setup_key = "netbird-setup-key"
    key = "setup-key"
}
tls_cert_secret = "nginx-tls-secret"

jellyfin: ClusterService = {
    name: "jellyfin"
    namespace: "jellyfin"
    cluster_url: "jellyfin.jellyfin.svc.cluster.local:8096"
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
}

vault: ClusterService = {
    name: "vault"
    namespace: "vault"
    cluster_url: "vault.vault.svc.cluster.local:8200"
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
}

flux_source_controller: ClusterService = {
    name: "flux-source-controller"
    namespace: "flux-system"
    cluster_url: "source-controller.flux-system.svc.cluster.local:80"
    gatus.endpointType: "tcp"
    gatus.conditions: ["[CONNECTED] == true"]
}

nginx_gateway: ClusterService = {
    name: "nginx-gateway"
    namespace: "nginx-gateway"
    cluster_url: "nginx-gateway.nginx-gateway.svc.cluster.local:80"
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
}

cert_manager: ClusterService = {
    name: "cert-manager"
    namespace: "cert-manager"
    cluster_url: "cert-manager.cert-manager.svc.cluster.local:9402"
    gatus.endpointType: "tcp"
    gatus.conditions: ["[CONNECTED] == true"]
}

cluster_services: [ClusterService] = [
    jellyfin,
    vault,
    flux_source_controller,
    nginx_gateway,
    cert_manager,
]