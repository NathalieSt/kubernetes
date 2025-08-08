import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared

out_dir = shared.cluster_path + "apps/wikijs/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
strict_tls_file = "strict-tls.yaml"
# Configs
service_name = "wikijs"
chart_url = "https://charts.js.wiki"
chart_name = "wiki"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
    strict_tls_file
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

strict_tls = {
    apiVersion: "security.istio.io/v1"
    kind: "PeerAuthentication"
    metadata: {
        name: "default"
        namespace: service_name
    }
    spec: {
        mtls: {
            mode: "STRICT"
        }
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
            postgresql: {
                enabled: False
                postgresqlHost: shared.postgres.cluster_url
                postgresqlPort: shared.postgres.port
                postgresqlDatabase: "wikijs"
                existingSecret: "postgres-creds-secret"
                existingSecretKey: "password"
            }
        }
    }
}
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(strict_tls, out_dir + strict_tls_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
