import ..lib.shared

cluster_path = "../../../../cluster/"
global_vault_auth_namespace = "vault-secrets-operator"
ca_issuer = "my-ca-issuer"
postgres_creds_secret_name = "postgres-creds-secret"
tls_cert_secret = "nginx-tls-secret"
smb_storage_class_name = "smb"
nfs_storage_class_name = "nfs-client"

jellyfin: shared.ClusterService = {
    name: "jellyfin"
    namespace: "jellyfin"
    cluster_url: "jellyfin.jellyfin.svc.cluster.local"
    port: 8096
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    gateway_config: {
        dns_name: "jellyfin.cluster"
        requires_websockets: True
        requires_larger_body: True
    }
}

vault: shared.ClusterService = {
    name: "vault"
    namespace: "vault"
    cluster_url: "vault.vault.svc.cluster.local"
    port: 8200
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    gateway_config: {
        dns_name: "vault.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

postgres: shared.ClusterService = {
    name: "postgres"
    namespace: "postgres"
    cluster_url: "postgres-postgresql.postgres.svc.cluster.local"
    port: 5432
    gatus.endpointType: "tcp"
    gatus.conditions: ["[CONNECTED] == true"]
}

wikijs: shared.ClusterService = {
    name: "wikijs"
    namespace: "wikijs"
    cluster_url: "wikijs.wikijs.svc.cluster.local"
    port: 80
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    gateway_config: {
        dns_name: "wikijs.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

gluetun_proxy: shared.ClusterService = {
    name: "gluetun-proxy"
    namespace: "gluetun-proxy"
    cluster_url: "gluetun-proxy.gluetun-proxy.svc.cluster.local"
    port: 8888
    gatus.endpointType: "tcp"
    gatus.conditions: ["[CONNECTED] == true"]
}

valkey: shared.ClusterService = {
    name: "valkey"
    namespace: "valkey"
    cluster_url: "valkey.valkey.svc.cluster.local"
    port: 6379
    gatus.endpointType: "tcp"
    gatus.conditions: ["[CONNECTED] == true"]
}

searxng: shared.ClusterService = {
    name: "searxng"
    namespace: "searxng"
    cluster_url: "searxng.searxng.svc.cluster.local"
    port: 8080
    gatus.endpointType: "http"
    gatus.conditions: ["[STATUS] == 200"]
    gateway_config: {
        dns_name: "searxng.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

nginx_gateway: shared.ClusterService = {
    name: "nginx-gateway"
    namespace: "searxng"
    cluster_url: "nginx-gateway.nginx-gateway.svc.cluster.local"
    port: 80
    gatus.conditions: ["[CONNECTED] == true"]
    gatus.endpointType: "http"
}

gatus: shared.ClusterService = {
    name: "gatus"
    namespace: "gatus"
    cluster_url: "gatus.gatus.svc.cluster.local"
    port: 80
    gatus.conditions: ["[CONNECTED] == true"]
    gatus.endpointType: "http"
    gateway_config: {
        dns_name: "gatus.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

redis = shared.ClusterService {
    name: "redis"
    namespace: "redis"
    cluster_url: "redis.redis.svc.cluster.local"
    port: 6379
    gatus.conditions: ["[CONNECTED] == true"]
    gatus.endpointType: "http"
}

#FIXME: make gatus configs optional
dawarich: shared.ClusterService = {
    name: "dawarich"
    namespace: "dawarich"
    cluster_url: "dawarich.dawarich.svc.cluster.local"
    port: 3000
    gatus.conditions: ["[CONNECTED] == true"]
    gatus.endpointType: "http"
    gateway_config: {
        dns_name: "dawarich.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

cluster_services: [shared.ClusterService] = [
    jellyfin
    vault
    gatus
    postgres
    wikijs
    gluetun_proxy
    valkey
    searxng
    dawarich
]

exposed_services = filter service in cluster_services {
    service.gateway_config != Undefined
}
