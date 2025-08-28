import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import ...versions

out_dir = shared.cluster_path + "infrastructure/authelia/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
# Configs
service_name = shared.authelia.name
chart_url = "https://charts.authelia.com"
chart_name = "authelia"
chart_version = versions.helm_charts.authelia
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
        valuesFrom: [
            {
                kind: "Secret"
                # FIXME: in shared, replace "postgres_creds_secret_name" with "postgres_creds_secret" and put "name: str, key: str" inside 
                name: shared.postgres_creds_secret_name
                optional: False
                targetPath: "configMap.storage.postgres.password.value"
                valuesKey: "password"
            }
        ]
        values: {
            configMap.session.cookies: [
                "https://${shared.authelia.gateway_config.dns_name}.netbird.selfhosted"
            ]
            configMap.storage.postgres: {
                enabled: True
                address: "tcp://${shared.postgres.cluster_url}:${shared.postgres.port}"
                database: "authelia"
                password.secret_name: "authelia-pg-password"
                username: "postgres"
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
