import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ....lib.helpers
import ...shared
import ...versions

out_dir = shared.cluster_path + "apps/forgejo/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
# Configs
service_name = shared.forgejo.name
# FIXME: TODO:
chart_url = "oci://code.forgejo.org/forgejo-helm/forgejo"
chart_name = "forgejo"
chart_version = versions.helm_charts.forgejo
netbird_secret_name = shared.netbird_secret_name
netbird_secret_key = shared.netbird_secret_name
forgejo_postgres_db_name = "forgejo"
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
        ref.tag: "14.0.0"
    }
}
release: helm_flux.HelmRelease = {
    metadata.name = service_name
    spec: {
        releaseName: service_name
        chartRef: {
            kind: "OCIRepository"
            name: service_name
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
                targetPath: "gitea.config.database.PASSWD"
                valuesKey: "password"
            }
        ]
        values: {
            gitea: {
                config: {
                    database: {
                        # supported values are mysql, postgres, mssql, sqlite3
                        DB_TYPE: "postgres"
                        HOST: "${shared.postgres.cluster_url}:${shared.postgres.port}"
                        NAME: forgejo_postgres_db_name
                        USER: "postgres"
                    }
                    server: {
                        ROOT_URL: "https://${shared.forgejo.gateway_config.dns_name}.netbird.selfhosted"
                    }
                }
                queue: {
                    TYPE: "redis"
                    CONN_STR: "valkey://${shared.valkey.cluster_url}:${shared.valkey.port}/0?"
                }
                cache: {
                    ADAPTER: "redis"
                    HOST: "valkey://${shared.valkey.cluster_url}:${shared.valkey.port}/1"
                }
                session: {
                    PROVIDER: "redis"
                    PROVIDER_CONFIG: "valkey://${shared.valkey.cluster_url}:${shared.valkey.port}/2"
                }
            }
            persistence: {
                enabled: True
                storageClass: shared.nfs_remote_storage_class_name
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
