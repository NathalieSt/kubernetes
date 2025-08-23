import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import ...versions
import ..istio_shared

out_dir = shared.cluster_path + "istio/ingress/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repo_file = "repo.yaml"
# Configs
chart_name = "gateway"
chart_version = versions.helm_charts.istio_ingress
service_name = "istio-ingress"

namespace = k8s_wrapper.create_namespace(service_name, True)

kustomization = k8s_wrapper.create_kustomization(istio_shared.repo_name, [
    release_file
    namespace_file
    repo_file
])

# Helm configuration
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
                    name: istio_shared.repo_name
                }
            }
        }
        install.remediation.retries: 3
        timeout: "5m"
        interval: "24h"
        values: {
            autoscaling: {
                maxReplicas: 2
            }
            service: {
                loadBalancerIP: shared.istio_gateway_ip
            }
        }
    }
}

repo: helm_flux.HelmRepository = {
    metadata.name = istio_shared.repo_name
    kind: "HelmRepository"
    spec: {
        interval: "24h"
        url: istio_shared.chart_url
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(release, out_dir + release_file)
yaml.dump_to_file(repo, out_dir + repo_file)
