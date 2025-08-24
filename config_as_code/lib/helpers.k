import .flux

resources = {
    requests: {
        cpu: "50m"
        memory: "64Mi"
    }
    limits: {
        cpu: "100m"
        memory: "128Mi"
    }
}

create_netbird_sidecar = lambda setup_key_secret_name: str, setup_key_secret_key: str, service_name: str {
    sidecar = {
        name: "netbird-agent"
        image: "netbirdio/netbird:latest"
        env: [
            {
                name: "NB_SETUP_KEY"
                valueFrom: {
                    secretKeyRef: {
                        name: setup_key_secret_name
                        key: setup_key_secret_key
                    }
                }
            }
            {
                name: "NB_HOSTNAME"
                value: service_name
            }
            {
                name: "NB_MANAGEMENT_URL"
                value: "https://netbird.nathalie-stiefsohn.eu"
            }
        ]
        resources: resources
        securityContext: {
            privileged: True
        }
    }
    sidecar
}

create_cluster_kustomization = lambda service_name: str, kustomization_path: str, target_namespace: str, dependencies: [flux.FluxKustomizationDependency] {
    kustomization: flux.FluxKustomization = {
        apiVersion: "kustomize.toolkit.fluxcd.io/v1"
        kind: "Kustomization"
        metadata: {
            name: service_name
            namespace: "flux-system"
        }
        spec: {
            interval: "24h"
            dependsOn: dependencies
            path: kustomization_path
            prune: True
            sourceRef: {
                kind: "GitRepository"
                name: "flux-system"
            }
            targetNamespace: target_namespace
            timeout: "10m"
            wait: True
        }
    }
    kustomization
}
