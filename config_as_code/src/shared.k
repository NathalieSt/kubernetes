schema Netbird:
    setup_key: str
    key: str

schema ClusterService:
    name: str
    namespace: str
    cluster_url: str

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
}

vault: ClusterService = {
    name: "vault"
    namespace: "vault"
    cluster_url: "vault.vault.svc.cluster.local:8200"
}

flux_source_controller: ClusterService = {
    name: "flux-source-controller"
    namespace: "flux-system"
    cluster_url: "source-controller.flux-system.svc.cluster.local:80"
}

nginx_gateway: ClusterService = {
    name: "nginx-gateway"
    namespace: "nginx-gateway"
    cluster_url: "nginx-gateway.nginx-gateway.svc.cluster.local:80"
}

cluster_services: [ClusterService] = [
    jellyfin,
    vault,
    flux_source_controller,
    nginx_gateway,
]