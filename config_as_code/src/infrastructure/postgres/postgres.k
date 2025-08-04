import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "infrastructure/postgres/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
init_configmap_file = "init-configmap.yaml"
# Configs
service_name = "postgres"
chart_url = "oci://registry-1.docker.io/bitnamicharts/postgresql"
chart_name = "postgres"
init_configmap_name = "init-configmap"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
    init_configmap_file
])
namespace = k8s_wrapper.create_namespace(service_name)
# init configmap
init_configmap: core_api.ConfigMap = {
    metadata: {
        name: init_configmap_name
    }
    data: {
        "init-wikijs.sql": """
            SELECT 'CREATE DATABASE wikijs'
            WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'wikijs')\gexec
        
        """
    }
}
# Helm configuration
repo: helm_flux.HelmRepository = {
    metadata.name = service_name
    kind: "OCIRepository"
    spec: {
        interval: "24h"
        url: chart_url
        ref: {
            tag: "16.7.15"
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
            primary.persistence.storageClass: "nfs-client"
            primary.initdb.scriptsConfigMap: init_configmap_name
            backup.enabled: True
            backup.cronjob.storage.storageClass: "smb"
        }
    }
}
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(init_configmap, out_dir + init_configmap_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
