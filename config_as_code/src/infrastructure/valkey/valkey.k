import ....lib.k8s_wrapper
import ...shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "infrastructure/valkey/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
pvc_file = "pvc.yaml"
deployment_file = "deployment.yaml"
service_file = "service.yaml"
# Configs
service_name = "valkey"
label_app_name = "valkey"
k8s_service_name = "valkey"
pvc_name = "valkey"
port = 6379
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    pvc_file
    deployment_file
    service_file
])
namespace = k8s_wrapper.create_namespace(service_name)

pvc: core_api.PersistentVolumeClaim = {
    apiVersion: "v1"
    kind: "PersistentVolumeClaim"
    metadata: {
        name: "valkey-pvc"
        annotations: {
            "app.kubernetes.io/name": label_app_name
        }
    }
    spec: {
        storageClassName: "nfs-client"
        # FIXME: Is this the correct access mode???
        accessModes: [
            "ReadWriteMany"
        ]
        resources: {
            requests: {
                "storage": "10Gi"
            }
        }
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: "valkey-deployment"
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
                        name: "valkey"
                        image: "valkey/valkey:8-alpine3.22"
                        ports: [
                            {
                                name: "valkey"
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                name: "valkey-volume"
                                mountPath: "/data"
                            }
                        ]
                    }
                ]
                volumes: [
                    {
                        name: "valkey-volume"
                        persistentVolumeClaim: {
                            claimName: "valkey-pvc"
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
        }
    }
    spec: {
        ports: [
            {
                name: "valkey"
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
yaml.dump_to_file(service, out_dir + service_file)
