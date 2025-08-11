import ..lib.shared
import .istio.istio_shared

cluster_path = "../../../../cluster/"
global_vault_auth_namespace = "vault-secrets-operator"
ca_issuer = "my-ca-issuer"
postgres_creds_secret_name = "postgres-creds-secret"
tls_cert_secret = "nginx-tls-secret"
smb_storage_class_name = "smb"
nfs_storage_class_name = "nfs-client"
nfs_remote_storage_class_name = "nfs-remote"

jellyfin: shared.ClusterService = {
    name: "jellyfin"
    namespace: "jellyfin"
    cluster_url: "jellyfin.jellyfin.svc.cluster.local"
    port: 8096
    gatus_config.endpointType: "http"
    gatus_config.conditions: ["[STATUS] == 200"]
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
    gatus_config.endpointType: "http"
    gatus_config.conditions: ["[STATUS] == 200"]
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
}

gluetun_proxy: shared.ClusterService = {
    name: "gluetun-proxy"
    namespace: "gluetun-proxy"
    cluster_url: "gluetun-proxy.gluetun-proxy.svc.cluster.local"
    port: 8888
    gatus_config.endpointType: "tcp"
    gatus_config.conditions: ["[CONNECTED] == true"]
}

valkey: shared.ClusterService = {
    name: "valkey"
    namespace: "valkey"
    cluster_url: "valkey.valkey.svc.cluster.local"
    port: 6379
}

searxng: shared.ClusterService = {
    name: "searxng"
    namespace: "searxng"
    cluster_url: "searxng.searxng.svc.cluster.local"
    port: 8080
    gatus_config.endpointType: "http"
    gatus_config.conditions: ["[STATUS] == 200"]
    gateway_config: {
        dns_name: "searxng.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

gatus: shared.ClusterService = {
    name: "gatus"
    namespace: "gatus"
    cluster_url: "gatus.gatus.svc.cluster.local"
    port: 80
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
}

dawarich: shared.ClusterService = {
    name: "dawarich"
    namespace: "dawarich"
    cluster_url: "dawarich.dawarich.svc.cluster.local"
    port: 3000
    gatus_config.conditions: ["[CONNECTED] == true"]
    gatus_config.endpointType: "tcp"
    gateway_config: {
        dns_name: "dawarich.cluster"
        requires_websockets: True
        requires_larger_body: False
    }
    virtual_service_config: {
        # https protocol needs to be forwarded to dawarich to prevent
        # a mismatch between configured server url and origin header url
        client_protocol: "https"
    }
}

kiali: shared.ClusterService = {
    name: "kiali"
    namespace: "kiali-operator"
    cluster_url: "kiali.${istio_shared.namespace}.svc.cluster.local"
    port: 20001
    gateway_config: {
        dns_name: "kiali.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

heimdall: shared.ClusterService = {
    name: "heimdall"
    namespace: "heimdall"
    cluster_url: "heimdall.heimdall.svc.cluster.local"
    port: 80
    gateway_config: {
        dns_name: "heimdall.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

tandoor: shared.ClusterService = {
    name: "tandoor"
    namespace: "tandoor"
    cluster_url: "tandoor.tandoor.svc.cluster.local"
    port: 80
    gateway_config: {
        dns_name: "tandoor.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

istio_gateway_ip = "10.43.164.134"

cluster_services: [shared.ClusterService] = [
    jellyfin
    vault
    gatus
    postgres
    gluetun_proxy
    valkey
    searxng
    dawarich
    kiali
    heimdall
    tandoor
]

cluster_services_watched_by_gatus = filter service in cluster_services {
    service.gatus_config != Undefined
}

exposed_services = filter service in cluster_services {
    service.gateway_config != Undefined
}
