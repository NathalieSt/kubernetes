import k8s.api.core.v1 as core_api

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import ...versions
import yaml

out_dir = shared.cluster_path + "infrastructure/cnpg/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
# Configs
service_name = shared.cnpg.name
chart_url = "https://cloudnative-pg.github.io/charts"
chart_name = "cloudnative-pg"
chart_version = versions.helm_charts.cnpg
# k8s configuration
kustomization: k8s_wrapper.Kustomization = {
    metadata.name = service_name
    resources = [
        namespace_file
        release_file
        repository_file
    ]
}
namespace: core_api.Namespace = k8s_wrapper.create_namespace(service_name, True)

# Helm configuration
repo: helm_flux.HelmRepository = {
    metadata.name = service_name
    kind: "HelmRepository"
    spec: {
        interval: "24h"
        url: chart_url
    }
}
release: helm_flux.HelmRelease = {
    metadata.name = service_name
    spec: {
        releaseName: service_name
        chart: {
            spec: {
                chart: chart_name
                version: chart_version
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
        values: {}
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
