schema GatusConfig:
    endpointType: "http" | "tcp"
    conditions?: [str]

schema ClusterService:
    name: str
    namespace: str
    cluster_url: str
    port: int
    gatus: GatusConfig
    exposed: bool

schema ExternalService:
    name: str
    url: str

cluster_path = "../../../../cluster/"
global_vault_auth_namespace = "vault-secrets-operator"
ca_issuer = "my-ca-issuer"
postgres_creds_secret_name = "postgres-creds-secret"
tls_cert_secret = "nginx-tls-secret"
smb_storage_class_name = "smb"

jellyfin: ClusterService = {
    name: "jellyfin"
    namespace: "jellyfin"
    cluster_url: "jellyfin.jellyfin.svc.cluster.local"
    port: 8096
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    exposed: True
}

vault: ClusterService = {
    name: "vault"
    namespace: "vault"
    cluster_url: "vault.vault.svc.cluster.local"
    port: 8200
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    exposed: True
}

postgres: ClusterService = {
    name: "postgres"
    namespace: "postgres"
    cluster_url: "postgres-postgresql.postgres.svc.cluster.local"
    port: 5432
    gatus.endpointType: "tcp"
    gatus.conditions: ["[CONNECTED] == true"]
    exposed: False
}

wikijs: ClusterService = {
    name: "wikijs"
    namespace: "wikijs"
    cluster_url: "wikijs.wikijs.svc.cluster.local"
    port: 80
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    exposed: True
}

gluetun_proxy: ClusterService = {
    name: "gluetun-proxy"
    namespace: "gluetun-proxy"
    cluster_url: "gluetun-proxy.gluetun-proxy.svc.cluster.local"
    port: 8888
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    exposed: False
}

valkey: ClusterService = {
    name: "valkey"
    namespace: "valkey"
    cluster_url: "valkey.valkey.svc.cluster.local"
    port: 6379
    gatus.endpointType: "tcp"
    gatus.conditions: ["[CONNECTED] == true"]
    exposed: False
}

searxng: ClusterService = {
    name: "searxng"
    namespace: "searxng"
    cluster_url: "searxng.searxng.svc.cluster.local"
    port: 8080
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    exposed: True
}

cluster_services: [ClusterService] = [
    jellyfin
    vault
    postgres
    wikijs
    gluetun_proxy
    valkey
    searxng
]
