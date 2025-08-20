import ....lib.k8s_wrapper
import ....lib.helm_flux
import ....lib.istio.destination_rule
import ....lib.istio.peer_authenthication
import ...shared
import ...versions
import yaml

out_dir = shared.cluster_path + "infrastructure/vault/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
vso_peer_auth_file = "peer-auth.yaml"
vso_destination_rule_file = "destination-rule.yaml"
# Configs
service_name = shared.vault.name
chart_url = "https://helm.releases.hashicorp.com"
chart_name = "vault"
chart_version = versions.helm_charts.vault
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    release_file
    repository_file
    vso_peer_auth_file
    vso_destination_rule_file
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
            ui.enabled: True
            server.dataStorage.storageClass: shared.nfs_local_storage_class_v2_name
        }
    }
}

vso_destination_rule: destination_rule.DestinationRule = {
    metadata: {
        name: "vault-dr"
    }
    spec: {
        host: shared.vault.cluster_url
        trafficPolicy: {
            tls: {
                mode: "DISABLE"
            }
        }
    }
}

vso_peer_auth: peer_authenthication.PeerAuthenthication = {
    metadata: {
        name: "vault-permissive-mtls"
    }
    spec: {
        selector: {
            matchLabels: {
                "app.kubernetes.io/name": "vault"
            }
        }
        mtls: {
            mode: "PERMISSIVE"
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
yaml.dump_to_file(vso_destination_rule, out_dir + vso_destination_rule_file)
yaml.dump_to_file(vso_peer_auth, out_dir + vso_peer_auth_file)
