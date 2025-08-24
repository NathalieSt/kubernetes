import ..lib.shared
import .istio.istio_shared

cluster_path = "../../../../cluster/"
global_vault_auth_namespace = "vault-secrets-operator"
postgres_creds_secret_name = "postgres-creds-secret"
nfs_local_storage_class_v2_name = "nfs-local-v2"
nfs_remote_storage_class_name = "nfs-remote"

jellyfin: shared.ClusterService = {
    name: "jellyfin"
    namespace: "jellyfin"
    type: "app"
    cluster_url: "jellyfin.jellyfin.svc.cluster.local"
    port: 8096
    gateway_config: {
        dns_name: "jellyfin.cluster"
        requires_websockets: True
        requires_larger_body: True
    }
}

vault: shared.ClusterService = {
    name: "vault"
    namespace: "vault"
    type: "infrastructure"
    cluster_url: "vault.vault.svc.cluster.local"
    port: 8200
    gateway_config: {
        dns_name: "vault.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

postgres: shared.ClusterService = {
    name: "postgres"
    namespace: "postgres"
    type: "infrastructure"
    cluster_url: "postgres-postgresql.postgres.svc.cluster.local"
    port: 5432
}

gluetun_proxy: shared.ClusterService = {
    name: "gluetun-proxy"
    namespace: "gluetun-proxy"
    type: "infrastructure"
    cluster_url: "gluetun-proxy.gluetun-proxy.svc.cluster.local"
    port: 8888
}

valkey: shared.ClusterService = {
    name: "valkey"
    namespace: "valkey"
    type: "infrastructure"
    cluster_url: "valkey.valkey.svc.cluster.local"
    port: 6379
}

searxng: shared.ClusterService = {
    name: "searxng"
    namespace: "searxng"
    type: "infrastructure"
    cluster_url: "searxng.searxng.svc.cluster.local"
    port: 8080
    gateway_config: {
        dns_name: "searxng.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

redis = shared.ClusterService {
    name: "redis"
    namespace: "redis"
    type: "infrastructure"
    cluster_url: "redis.redis.svc.cluster.local"
    port: 6379
}

dawarich: shared.ClusterService = {
    name: "dawarich"
    namespace: "dawarich"
    type: "app"
    cluster_url: "dawarich.dawarich.svc.cluster.local"
    port: 3000
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
    type: "monitoring"
    cluster_url: "kiali.${istio_shared.namespace}.svc.cluster.local"
    port: 20001
    gateway_config: {
        dns_name: "kiali.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

mealie: shared.ClusterService = {
    name: "mealie"
    namespace: "mealie"
    type: "app"
    cluster_url: "mealie.mealie.svc.cluster.local"
    port: 9000
    gateway_config: {
        dns_name: "mealie.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

glance: shared.ClusterService = {
    name: "glance"
    namespace: "glance"
    type: "app"
    cluster_url: "glance.glance.svc.cluster.local"
    port: 8080
    gateway_config: {
        dns_name: "glance.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

istio_gateway_ip = "10.43.164.134"

cluster_services: [shared.ClusterService] = [
    jellyfin
    vault
    postgres
    gluetun_proxy
    valkey
    searxng
    dawarich
    kiali
    mealie
    glance
]

exposed_services = filter service in cluster_services {
    service.gateway_config != Undefined
}
