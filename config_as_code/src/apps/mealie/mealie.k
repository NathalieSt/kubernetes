import ....lib.k8s_wrapper
import ....lib.istio
import ...shared
import ...versions
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "apps/mealie/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
pvc_file = "pvc.yaml"
# Configs
service_name = shared.mealie.name
label_app_name = service_name
k8s_service_name = service_name
image_version = versions.docker_images.mealie
deployment_pvc_volume = "mealie-pvc"
port = shared.mealie.port
pvc_name = "mealie-pvc"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    pvc_file
    deployment_file
    k8s_service_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)

pvc: core_api.PersistentVolumeClaim = {
    metadata.name = pvc_name
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_local_storage_class_v2_name
        resources.requests.storage: "100Gi"
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: service_name
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        replicas: 1
        selector: {
            matchLabels: {
                "app.kubernetes.io/name": label_app_name
            }
        }
        template: {
            metadata: meta.ObjectMeta {
                labels: {
                    "app.kubernetes.io/name": label_app_name
                    "app.kubernetes.io/version": image_version
                }
            }
            spec = core_api.PodSpec {
                containers: [
                    {
                        name: service_name
                        image: "ghcr.io/mealie-recipes/mealie:${image_version}"
                        env: [
                            {name: "POSTGRES_SERVER", value: shared.postgres.cluster_url}
                            {name: "POSTGRES_PORT", value: "${shared.postgres.port}"}
                            {name: "POSTGRES_USERNAME", value: "postgres"}
                            {name: "POSTGRES_PASSWORD", valueFrom: {
                                secretKeyRef: {
                                    name: shared.postgres_creds_secret_name
                                    key: "password"
                                }
                            }}
                            {name: "POSTGRES_DB", value: "mealie"}
                            {name: "BASE_URL", value: "https://${shared.mealie.gateway_config.dns_name}.netbird.selfhosted"}
                        ]
                        ports: [
                            {
                                name: service_name
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                mountPath: "/app/data/"
                                name: deployment_pvc_volume
                            }
                        ]
                    }
                ]
                volumes: [
                    {
                        name: deployment_pvc_volume
                        persistentVolumeClaim: {
                            claimName: pvc_name
                        }
                    }
                ]
            }
        }
    }
}

service: core_api.Service = {
    apiVersion: "v1"
    kind: "Service"
    metadata: {
        name: k8s_service_name
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        ports: [
            {
                name: "mealie"
                port: port
                targetPort: port
            }
        ]
        selector: {
            "app.kubernetes.io/name": label_app_name
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(pvc, out_dir + pvc_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(service, out_dir + k8s_service_file)
