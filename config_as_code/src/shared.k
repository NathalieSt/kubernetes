import ..lib.shared
import .istio.istio_shared

cluster_path = "../../../../cluster/"
global_vault_auth_namespace = "vault-secrets-operator"
postgres_creds_secret_name = "postgres-creds-secret"
netbird_secret_name = "netbird-setup-key-vault-secret"
netbird_secret_key = "setup-key"
nfs_local_storage_class_v2_name = "nfs-local-v2"
nfs_remote_storage_class_name = "nfs-remote"
timezone = "Europe/Vienna"

# Cluster Services
## Apps
kustomization_apps_path = "./cluster/apps"

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
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_apps_path}/jellyfin"
    }
    keda_scaling: {
        desiredReplicas: "1"
        start: "0 7 * * *"
        end: "0 23 * * *"
        timezone: timezone
    }
}

searxng: shared.ClusterService = {
    name: "searxng"
    namespace: "searxng"
    type: "app"
    cluster_url: "searxng.searxng.svc.cluster.local"
    port: 8080
    gateway_config: {
        dns_name: "searxng.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
    flux_kustomization: {
        dependencies: [{name: valkey.name}]
        kustomization_path: "${kustomization_apps_path}/searxng"
    }
    keda_scaling: {
        desiredReplicas: "1"
        start: "0 7 * * *"
        end: "0 23 * * *"
        timezone: timezone
    }
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
    flux_kustomization: {
        dependencies: [
            {name: reflector.name}
            {name: redis.name}
            {name: postgres.name}
        ]
        kustomization_path: "${kustomization_apps_path}/dawarich"
    }
    keda_scaling: {
        desiredReplicas: "1"
        start: "0 8 * * *"
        end: "0 22 * * *"
        timezone: timezone
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
    flux_kustomization: {
        dependencies: [
            {name: reflector.name}
            {name: postgres.name}
        ]
        kustomization_path: "${kustomization_apps_path}/mealie"
    }
    keda_scaling: {
        desiredReplicas: "1"
        start: "0 9 * * *"
        end: "0 21 * * *"
        timezone: timezone
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
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_apps_path}/glance"
    }
    keda_scaling: {
        desiredReplicas: "1"
        start: "0 7 * * *"
        end: "0 23 * * *"
        timezone: timezone
    }
}

forgejo: shared.ClusterService = {
    name: "forgejo"
    namespace: "forgejo"
    type: "app"
    cluster_url: "forgejo-http.forgejo.svc.cluster.local"
    port: 3000
    gateway_config: {
        dns_name: "code.cluster"
        requires_websockets: False
        requires_larger_body: False
    }
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_apps_path}/forgejo"
    }
    keda_scaling: {
        desiredReplicas: "1"
        start: "0 7 * * *"
        end: "0 2 * * *"
        timezone: timezone
    }
}

## Infrastructure
kustomization_infrastructure_path = "./cluster/infrastructure"

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
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_infrastructure_path}/vault"
    }
}

vault_secrets_operator: shared.ClusterService = {
    name: "vault-secrets-operator"
    namespace: "vault-secrets-operator"
    type: "infrastructure"
    flux_kustomization: {
        dependencies: [{name: vault.name}]
        kustomization_path: "${kustomization_infrastructure_path}/vault-secrets-operator"
    }
}

postgres: shared.ClusterService = {
    name: "postgres"
    namespace: "postgres"
    type: "infrastructure"
    cluster_url: "postgres-rw.postgres.svc.cluster.local"
    port: 5432
    flux_kustomization: {
        dependencies: [{name: reflector.name}]
        kustomization_path: "${kustomization_infrastructure_path}/postgres"
    }
}

gluetun_proxy: shared.ClusterService = {
    name: "gluetun-proxy"
    namespace: "gluetun-proxy"
    type: "infrastructure"
    cluster_url: "gluetun-proxy.gluetun-proxy.svc.cluster.local"
    port: 8888
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_infrastructure_path}/gluetun-proxy"
    }
}

valkey: shared.ClusterService = {
    name: "valkey"
    namespace: "valkey"
    type: "infrastructure"
    cluster_url: "valkey.valkey.svc.cluster.local"
    port: 6379
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_infrastructure_path}/valkey"
    }
}

redis: shared.ClusterService = {
    name: "redis"
    namespace: "redis"
    type: "infrastructure"
    cluster_url: "redis.redis.svc.cluster.local"
    port: 6379
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_infrastructure_path}/redis"
    }
}

reflector: shared.ClusterService = {
    name: "reflector"
    namespace: "reflector"
    type: "infrastructure"
    flux_kustomization: {
        dependencies: [{name: vault_secrets_operator.name}]
        kustomization_path: "${kustomization_infrastructure_path}/reflector"
    }
}

caddy: shared.ClusterService = {
    name: "caddy"
    namespace: "caddy"
    type: "infrastructure"
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_infrastructure_path}/caddy"
    }
}

csi_driver_nfs: shared.ClusterService = {
    name: "csi-driver-nfs"
    namespace: "csi-driver-nfs"
    type: "infrastructure"
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_infrastructure_path}/csi-driver-nfs"
    }
}

cnpg: shared.ClusterService = {
    name: "cnpg"
    namespace: "cnpg-system"
    type: "infrastructure"
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_infrastructure_path}/cnpg"
    }
}

keda: shared.ClusterService = {
    name: "keda"
    namespace: "keda"
    type: "infrastructure"
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_infrastructure_path}/keda"
    }
}

## Istio
kustomization_istio_path = "./cluster/istio"
istio_system_namespace = "istio-system"

istio_repo: shared.ClusterService = {
    name: "istio-repo"
    namespace: istio_system_namespace
    type: "istio"
    flux_kustomization: {
        dependencies: []
        kustomization_path: "${kustomization_istio_path}/repo"
    }
}

istio_base: shared.ClusterService = {
    name: "istio-base"
    namespace: istio_system_namespace
    type: "istio"
    flux_kustomization: {
        dependencies: [{name: istio_repo.name}]
        kustomization_path: "${kustomization_istio_path}/base"
    }
}

istiod: shared.ClusterService = {
    name: "istiod"
    namespace: istio_system_namespace
    type: "istio"
    flux_kustomization: {
        dependencies: [{name: istio_base.name}]
        kustomization_path: "${kustomization_istio_path}/istiod"
    }
}

istio_prometheus: shared.ClusterService = {
    name: "prometheus"
    namespace: istio_system_namespace
    type: "istio"
    flux_kustomization: {
        dependencies: [{name: istiod.name}]
        kustomization_path: "${kustomization_istio_path}/prometheus"
    }
}

istio_ingress: shared.ClusterService = {
    name: "istio-ingress"
    namespace: "istio-ingress"
    type: "istio"
    flux_kustomization: {
        dependencies: [{name: istiod.name}]
        kustomization_path: "${kustomization_istio_path}/ingress"
    }
}

istio_security: shared.ClusterService = {
    name: "istio-security"
    namespace: istio_system_namespace
    type: "istio"
    flux_kustomization: {
        dependencies: [
            {name: istiod.name}
            {name: istio_ingress.name}
        ]
        kustomization_path: "${kustomization_istio_path}/security"
    }
}

istio_networking: shared.ClusterService = {
    name: "istio-networking"
    namespace: istio_system_namespace
    type: "istio"
    flux_kustomization: {
        dependencies: [
            {name: istiod.name}
            {name: istio_ingress.name}
        ]
        kustomization_path: "${kustomization_istio_path}/networking"
    }
}

## Monitoring
kustomization_monitoring_path = "./cluster/monitoring"

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
    flux_kustomization: {
        dependencies: [{name: istio_prometheus.name}]
        kustomization_path: "${kustomization_monitoring_path}/kiali"
    }
}

istio_gateway_ip = "10.43.164.134"

cluster_app_services: [shared.ClusterService] = [
    jellyfin
    forgejo
    dawarich
    searxng
    mealie
    glance
]

cluster_infrastructure_services: [shared.ClusterService] = [
    vault
    vault_secrets_operator
    reflector
    postgres
    gluetun_proxy
    valkey
    redis
    caddy
    csi_driver_nfs
    cnpg
    keda
]

cluster_istio_services: [shared.ClusterService] = [
    istio_repo
    istio_base
    istiod
    istio_prometheus
    istio_ingress
    istio_security
    istio_networking
]

cluster_monitoring_services: [shared.ClusterService] = [
    kiali
]

cluster_services: [shared.ClusterService] = cluster_app_services + cluster_infrastructure_services + cluster_monitoring_services

exposed_services = filter service in cluster_services {
    service.gateway_config != Undefined
}
