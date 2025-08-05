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

wikijs: shared.ClusterService = {
    name: "wikijs"
    namespace: "wikijs"
    cluster_url: "wikijs.wikijs.svc.cluster.local"
    port: 80
    gatus_config.endpointType: "http"
    gatus_config.conditions: ["[STATUS] == 200"]
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

nginx_gateway: shared.ClusterService = {
    name: "nginx-gateway"
    namespace: "searxng"
    cluster_url: "nginx-gateway.nginx-gateway.svc.cluster.local"
    port: 80
    gatus_config.conditions: ["[CONNECTED] == true"]
    gatus_config.endpointType: "http"
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
    gatus_config.endpointType: "http"
    gateway_config: {
        dns_name: "dawarich.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
}

registry: shared.ClusterService = {
    name: "registry"
    namespace: "registry"
    cluster_url: "registry.registry.svc.cluster.local"
    port: 5000
    gatus_config.conditions: ["[CONNECTED] == true"]
    gatus_config.endpointType: "http"
    gateway_config: {
        dns_name: "registry.cluster"
        requires_websockets: False
        requires_larger_body: True
        extra_location_config: """
            if ($http_user_agent ~ "^(docker\/1\.(3|4|5(?!\.[0-9]-dev))|Go ).*$" ) {
                return 404;
            }

            # To add basic authentication to v2 use auth_basic setting.
            auth_basic "Registry realm";
            auth_basic_user_file /etc/nginx/conf.d/nginx.htpasswd;

            ## If $docker_distribution_api_version is empty, the header is not added.
            ## See the map directive above where this variable is defined.
            add_header 'Docker-Distribution-Api-Version' $docker_distribution_api_version always;
        """
        extra_server_config: """
            # required to avoid HTTP 411: see Issue #1486 (https://github.com/moby/moby/issues/1486)
            chunked_transfer_encoding on;
        """
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
    registry
]

cluster_services_watched_by_gatus = filter service in cluster_services {
    service.gatus_config != Undefined
}

exposed_services = filter service in cluster_services {
    service.gateway_config != Undefined
}
