schema HelmCharts:
    gatus: str
    csi_driver_nfs: str
    vault: str
    vault_secrets_operator: str
    postgres: str
    reflector: str
    istio_base: str
    istio_ingress: str
    istiod: str
    kiali: str
    jellyfin: str
    forgejo: str
    cnpg: str

schema DockerImages:
    caddy: str
    gluetun: str
    redis: str
    valkey: str
    prometheus: str
    glance: str
    searxng: str
    mealie: str
    dawarich: str
    postgres: str

docker_images: DockerImages = {
    caddy: "2.10.0-alpine"
    gluetun: "v3.40"
    redis: "8.2.1"
    valkey: "8-alpine3.22"
    prometheus: "v3.5.0"
    glance: "v0.8.4"
    searxng: "2025.8.3-2e62eb5"
    mealie: "v3.0.2"
    dawarich: "0.30.10"
    postgres: "17.6"
}

helm_charts: HelmCharts = {
    gatus: "1.3.1"
    csi_driver_nfs: "4.11.0"
    vault: "0.30.1"
    vault_secrets_operator: "0.10.0"
    postgres: "16.7.15"
    reflector: "9.1.27"
    istio_base: "1.27.0"
    istio_ingress: "1.27.0"
    istiod: "1.27.0"
    kiali: "2.14.0"
    jellyfin: "2.3.0"
    forgejo: "TODO"
    cnpg: "0.26.0"
}
