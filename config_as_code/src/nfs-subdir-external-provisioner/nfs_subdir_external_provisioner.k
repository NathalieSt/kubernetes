import ...lib.k8s_wrapper
import ...lib.helm_flux
import ..shared
import yaml

out_dir = shared.cluster_path + "infrastructure/nfs-subdir-external-provisioner/"

# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"

# Configs
service_name = "nfs-subdir-external-provisioner"
chart_url = "https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner"
chart_name = "nfs-subdir-external-provisioner"

# k8s configuration
kustomization = k8s_wrapper.create_kustomization(
    service_name,
    [
        namespace_file,
        release_file,
        repository_file,
    ]
)
namespace = k8s_wrapper.create_namespace(service_name)

# Helm configuration
repo: helm_flux.HelmRepository = {
    metadata.name = service_name
    spec: {
        interval: "24h",
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
            image: {
                repository: "registry.k8s.io/sig-storage/nfs-subdir-external-provisioner"
                tag: "v4.0.2 # {\"$imagepolicy\": \"flux-system:nfs-subdir-external-provisioner-image-policy:tag\"}"
            }
            nfs: {
                server: "raspberry-pi-4-storage"
                path: "/mnt/external_ssd"
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