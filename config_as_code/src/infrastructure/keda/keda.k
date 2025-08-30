import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ....lib.keda
import ...shared
import ...versions

out_dir = shared.cluster_path + "infrastructure/keda/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
scaled_objects_file = "scaled-objects.yaml"
# Configs
service_name = shared.keda.name
chart_url = "https://kedacore.github.io/charts"
chart_name = "keda"
chart_version = versions.helm_charts.keda
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
    scaled_objects_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)
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
        values: {
            # -- Pod annotations for KEDA operator
            keda: {
                "traffic.sidecar.istio.io/excludeInboundPorts": "9666"
                "traffic.sidecar.istio.io/excludeOutboundPorts": "9443,6443"
            }
            # -- Pod annotations for KEDA Metrics Adapter
            metricsAdapter: {
                "traffic.sidecar.istio.io/excludeInboundPorts": "6443"
                "traffic.sidecar.istio.io/excludeOutboundPorts": "9666,9443"
            }
            # -- Pod annotations for KEDA Admission webhooks
            webhooks: {
                "traffic.sidecar.istio.io/excludeInboundPorts": "9443"
                "traffic.sidecar.istio.io/excludeOutboundPorts": "9666,6443"
            }
        }
    }
}

scaled_services = filter service in shared.cluster_services {
    service.keda_scaling != Undefined
}

scaled_objects: [keda.ScaledObject] = [{
    metadata: {
        name: "${scaled_service.name}-scaledobject"
        namespace: "${scaled_service.namespace}"
    }
    spec: {
        minReplicaCount: 0
        cooldownPeriod: 300
        scaleTargetRef: {
            name: "${scaled_service.name}"
        }
        triggers: [
            {
                type: "cron"
                metadata: {
                    desiredReplicas: "${scaled_service.keda_scaling.desiredReplicas}"
                    start: scaled_service.keda_scaling.start
                    end: scaled_service.keda_scaling.end
                    timezone: "Europe/Vienna"
                }
            }
        ]
    }
} for scaled_service in scaled_services]

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
yaml.dump_all_to_file(scaled_objects, out_dir + scaled_objects_file)
