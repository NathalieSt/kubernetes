import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ....lib.shared as lib_shared
import ....lib.helpers
import ...shared

out_dir = shared.cluster_path + "monitoring/gatus/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
# Configs
service_name = "gatus"
chart_url = "https://twin.github.io/helm-charts"
chart_name = "gatus"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
])
namespace = k8s_wrapper.create_namespace(service_name)
# Helm configuration
repo: helm_flux.HelmRepository = {
    metadata.name = service_name
    kind: "HelmRepository"
    spec: {
        interval: "24h"
        url: chart_url
    }
}

external_services: [lib_shared.ExternalService] = [
    {
        name: "netbird"
        url: "https://netbird.nathalie-stiefsohn.eu"
    }
    {
        name: "revolt"
        url: "https://chat.nathalie-stiefsohn.eu"
    }
]

cluster_internal_endpoints = [helpers.create_gatus_endpoint(service.name, service.gatus_config.endpointType + "://" + service.cluster_url + ":" + str(service.port), "GET", service.gatus_config.conditions, "60s", "Cluster Services") for service in shared.cluster_services_watched_by_gatus]

external_endpoints = [helpers.create_gatus_endpoint(service.name, service.url, "GET", ["[STATUS] == 200"], "60s", "External Services") for service in external_services]

release: helm_flux.HelmRelease = {
    metadata.name = service_name
    spec: {
        releaseName: service_name
        chart: {
            spec: {
                chart: chart_name
                version: "*"
                interval: "24h"
                sourceRef: {
                    kind: "HelmRepository"
                    name: service_name
                }
            }
        }
        install.remediation.retries: 3
        timeout: "5m"
        interval: "24h"
        values: {
            persistence.enabled: True
            config: {
                storage.type: "sqlite"
                storage.path: "/data/data.db"
                endpoints: cluster_internal_endpoints + external_endpoints
            }
        }
    }
}
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
