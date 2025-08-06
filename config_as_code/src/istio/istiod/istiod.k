import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import ..istio_shared

out_dir = shared.cluster_path + "istio/istiod/"
# Filenames
kustomization_file = "kustomization.yaml"
release_file = "release.yaml"
# Configs
service_name = "istiod"
chart_name = "istiod"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    release_file
])

# Helm configuration
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
                    name: istio_shared.repo_name
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
# Write helm configuration to files
yaml.dump_to_file(release, out_dir + release_file)
