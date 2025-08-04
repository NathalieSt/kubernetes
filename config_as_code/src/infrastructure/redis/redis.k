import ....lib.k8s_wrapper
import ...shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "infrastructure/redis/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
pvc_file = "redis-pvc.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
# Configs
service_name = shared.redis.name
label_app_name = "redis"
k8s_service_name = "redis"
pvc_name = "redis-pvc"
port = shared.redis.port
pvc_volume_name = "redis-volume"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    pvc_file
    deployment_file
    k8s_service_file
])
namespace = k8s_wrapper.create_namespace(service_name)

pvc: core_api.PersistentVolumeClaim = {
    metadata.name = pvc_name
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_storage_class_name
        resources.requests.storage: "10Gi"
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: service_name
        labels: {
            "app.kubernetes.io/name": label_app_name
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
                }
            }
            spec = core_api.PodSpec {
                containers: [
                    {
                        name: "redis"
                        image: "redis:latest"
                        ports: [
                            {
                                name: "redis"
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                mountPath: "/data"
                                name: pvc_volume_name
                            }
                        ]
                    }
                ]
                volumes: [
                    {
                        name: pvc_volume_name
                        persistentVolumeClaim: {
                            claimName: pvc_name
                        }
                    }
                ]
            }
        }
    }
}

k8s_service: core_api.Service = {
    apiVersion: "v1"
    kind: "Service"
    metadata: {
        name: k8s_service_name
        labels: {
            "app.kubernetes.io/name": label_app_name
        }
    }
    spec: {
        ports: [
            {
                name: "redis"
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
yaml.dump_to_file(k8s_service, out_dir + k8s_service_file)
