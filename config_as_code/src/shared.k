

schema GatusConfig:
    endpointType: "http" | "tcp"
    conditions?: [str]

schema ClusterService:
    name: str
    namespace: str
    cluster_url: str
    port: int
    gatus: GatusConfig

cluster_path = "../../../../cluster/"
global_vault_auth_namespace = "vault-secrets-operator"
ca_issuer = "my-ca-issuer"
postgres_creds_secret_name = "postgres-creds-secret"
tls_cert_secret = "nginx-tls-secret"
smb_storage_class_name= "smb"

jellyfin: ClusterService = {
    name: "jellyfin"
    namespace: "jellyfin"
    cluster_url: "jellyfin.jellyfin.svc.cluster.local"
    port: 8096
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
}

vault: ClusterService = {
    name: "vault"
    namespace: "vault"
    cluster_url: "vault.vault.svc.cluster.local"
    port: 8200
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
}

postgres: ClusterService = {
    name: "postgres"
    namespace: "postgres"
    cluster_url: "postgres-postgresql.postgres.svc.cluster.local"
    port: 5432
    gatus.endpointType: "tcp"
    gatus.conditions: ["[CONNECTED] == true"]
}

wikijs: ClusterService = {
    name: "wikijs"
    namespace: "wikijs"
    cluster_url: "wikijs.wikijs.svc.cluster.local"
    port: 80
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
} 

cluster_services: [ClusterService] = [
    jellyfin,
    vault,
    postgres,
    wikijs
]