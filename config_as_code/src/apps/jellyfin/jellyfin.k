import k8s.api.core.v1 as core_api
import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import ...versions
import ....lib.keda

out_dir = shared.cluster_path + "apps/jellyfin/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
pvc_file = "pvc.yaml"
scaled_object_file = "scaled-object.yaml"
# Configs
service_name = shared.jellyfin.name
namespace_name = shared.jellyfin.namespace
chart_url = "https://jellyfin.github.io/jellyfin-helm"
chart_name = "jellyfin"
chart_version = versions.helm_charts.jellyfin
pvc_name = "jellyfin-library-storage-pvc"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
    pvc_file
    scaled_object_file
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
            persistence.media.existingClaim: pvc_name
        }
    }
}

pvc: core_api.PersistentVolumeClaim = {
    metadata.name = pvc_name
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_local_storage_class_v2_name
        resources.requests.storage: "100Gi"
    }
}

scaled_object: keda.ScaledObject = {
    metadata: {
        name: "${service_name}-scaledobject"
        namespace: "${namespace_name}"
    }
    spec: {
        minReplicaCount: 0
        cooldownPeriod: 300
        scaleTargetRef: {
            name: "${service_name}"
        }
        triggers: [
            {
                type: "cron"
                metadata: shared.jellyfin.keda_scaling
            }
        ]
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
yaml.dump_to_file(pvc, out_dir + pvc_file)
yaml.dump_to_file(scaled_object, out_dir + scaled_object_file)
