import ...lib.k8s_wrapper
import ...lib.helm_flux
import ...lib.cert_manager
import ..shared
import yaml

out_dir = shared.cluster_path + "infrastructure/cert-manager/"

# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
cert_manager_file = "cluster-issuer.yaml"

# Configs
service_name = "cert-manager"
chart_url = "https://charts.jetstack.io"
chart_name = "cert-manager"
root_issuer_name = "selfsigned-issuer"
root_certificate_name = "root-cert"

# k8s configuration
kustomization = k8s_wrapper.create_kustomization(
    service_name,
    [
        namespace_file,
        release_file,
        repository_file,
        cert_manager_file,
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
                version: ""
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
            installCRDs: True
            replicaCount: 1
            resources.requests: {
                cpu: "10m"
                memory: "64Mi"
            }
            resources.limits: {
                cpu: "100m"
                memory: "128Mi"
            }
            webhook.replicaCount: 1
            webhook.resources.requests: {    
                cpu: "10m"
                memory: "32Mi"
            }
            webhook.resources.limits: {
                cpu: "100m"
                memory: "64Mi"
            }
            cainjector.replicaCount: 1  
            cainjector.resources.requests: {
                cpu: "10m"
                memory: "64Mi"
            }
            cainjector.resources.limits: {
                cpu: "100m"
                memory: "128Mi"
            }
        }
    }
}

# Cert Manager configuration
root_issuer: cert_manager.ClusterIssuer = {
    metadata.name: root_issuer_name
    spec: {
        selfSigned: {}
    }
}

root_certificate: cert_manager.Certificate = {
    metadata.name: "my-selfsigned-ca"
    spec: {
        isCa: True
        commonName: "my-selfsigned-ca"
        secretName: root_certificate_name
        privateKey: {
            algorithm: "ECDSA",
            size: 521
        }
        issuerRef: {
            name: root_issuer_name
            kind: "ClusterIssuer"
            group: "cert-manager.io"
        }
    }
}

cluster_issuer: cert_manager.ClusterIssuer = {
    metadata: {
        name: "my-ca-issuer"
    }
    spec: {
        ca: {
            secretName: root_certificate_name
        }
    }
}

cert_manger_config = [
    root_issuer,
    root_certificate,
    cluster_issuer,
]

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
# Write cert-manager configuration to files
yaml.dump_all_to_file(cert_manger_config, out_dir + cert_manager_file)