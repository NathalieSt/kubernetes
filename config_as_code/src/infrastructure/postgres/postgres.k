import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import ...versions
import yaml

out_dir = shared.cluster_path + "infrastructure/postgres/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
# Configs
service_name = "postgres"
chart_url = "oci://registry-1.docker.io/bitnamicharts/postgresql"
chart_name = "postgres"
chart_version = versions.helm_charts.postgres
init_configmap_name = "init-configmap"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)

# Helm configuration
repo: helm_flux.HelmRepository = {
    metadata.name = service_name
    kind: "OCIRepository"
    spec: {
        interval: "24h"
        url: chart_url
        ref: {
            tag: chart_version
        }
    }
}
release: helm_flux.HelmRelease = {
    metadata.name = service_name
    spec: {
        releaseName: service_name
        chartRef: {
            kind: "OCIRepository"
            name: "postgres"
        }
        install.remediation.retries: 3
        timeout: "5m"
        interval: "24h"
        valuesFrom: [
            {
                kind: "Secret"
                name: shared.postgres_creds_secret_name
                optional: False
                valuesKey: "password"
                targetPath: "auth.postgresPassword"
            }
        ]
        values: {
            primary.persistence.storageClass: shared.nfs_remote_storage_class_name
        }
    }
}
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
