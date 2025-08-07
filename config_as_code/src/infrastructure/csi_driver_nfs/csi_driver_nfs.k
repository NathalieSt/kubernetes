import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import yaml

out_dir = shared.cluster_path + "infrastructure/csi-driver-nfs/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
remote_storage_class_file = "remote-storage-class.yaml"
# Configs
service_name = "csi-driver-nfs"
chart_url = "https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/charts"
chart_name = "csi-driver-nfs"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
    remote_storage_class_file
])
namespace = k8s_wrapper.create_namespace(service_name)
# Helm configuration
repo: helm_flux.HelmRepository = {
    apiVersion: "source.toolkit.fluxcd.io/v1beta2"
    metadata.name = service_name
    kind: "HelmRepository"
    spec: {
        interval: "24h"
        url: chart_url
    }
}
release: helm_flux.HelmRelease = {
    apiVersion: "helm.toolkit.fluxcd.io/v2"
    kind: "HelmRelease"
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
        values: {}
    }
}

# the package seems to be broken
# to get IDE support, import "import k8s.api.storage.v1 as storage_api" and use storage_api for schemas
# remote storage class
#import k8s.api.storage.v1 as storage_api
remote_storage_class = {
    apiVersion: "storage.k8s.io/v1"
    kind: "StorageClass"
    metadata: {
        name: shared.nfs_remote_storage_class_name
    }
    provisioner: "nfs.csi.k8s.io"
    parameters: {
        "server": "remote-fs.netbird.selfhosted"
        "share": "/mnt/HC_Volume_103061115"
    }
    reclaimPolicy: "Delete"
    volumeBindingMode: "Immediate"
    allowVolumeExpansion: True
    mountOptions: [
        "nfsvers=4.1"
    ]
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
## Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
## Write storage classes to files
yaml.dump_to_file(remote_storage_class, out_dir + remote_storage_class_file)
