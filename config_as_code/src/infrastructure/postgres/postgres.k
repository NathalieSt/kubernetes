import ....lib.k8s_wrapper
import ....lib.cnpg
import ...shared
import ...versions
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "infrastructure/postgres/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
pvc_file = "pvc.yaml"
deployment_file = "deployment.yaml"
service_file = "service.yaml"
cnpg_cluster_file = "cnpg-cluster.yaml"
# Configs
service_name = shared.postgres.name
label_app_name = service_name
k8s_service_name = service_name
pvc_name = "${service_name}-pvc"
deployment_volume_name = "pvc-volume"
port = shared.postgres.port
image_version = versions.docker_images.postgres
cnpg_cluster_name = "postgres-v2"
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    #pvc_file
    #deployment_file
    #service_file
    cnpg_cluster_file
])
namespace = k8s_wrapper.create_namespace(shared.postgres.namespace, True)

pvc: core_api.PersistentVolumeClaim = {
    apiVersion: "v1"
    kind: "PersistentVolumeClaim"
    metadata: {
        name: pvc_name
        annotations: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        storageClassName: shared.nfs_remote_storage_class_name
        accessModes: [
            "ReadWriteMany"
        ]
        resources: {
            requests: {
                "storage": "20Gi"
            }
        }
    }
}

postgres_env: [core_api.EnvVar] = [
    {name: "POSTGRES_USER", value: "postgres"}
    {name: "POSTGRES_PASSWORD", valueFrom: {
        secretKeyRef: {
            name: shared.postgres_creds_secret_name
            key: "password"
        }
    }}
]

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: "${service_name}-deployment"
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
                        name: "${service_name}"
                        image: "postgres:${image_version}"
                        ports: [
                            {
                                name: "${service_name}"
                                containerPort: port
                            }
                        ]
                        env: postgres_env
                        volumeMounts: [
                            {
                                name: deployment_volume_name
                                mountPath: "/var/lib/postgresql/data"
                            }
                        ]
                    }
                ]
                volumes: [
                    {
                        name: deployment_volume_name
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
                name: "tcp-${service_name}"
                port: port
                targetPort: port
            }
        ]
        selector: {
            "app.kubernetes.io/name": label_app_name
        }
    }
}

cnpg_cluster: cnpg.Cluster = {
    metadata: {
        name: "postgres-v2"
    }
    spec: {
        instances: 3
        storage: {
            size: "20Gi"
            storageClass: shared.nfs_remote_storage_class_name
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)

#yaml.dump_to_file(pvc, out_dir + pvc_file)
#yaml.dump_to_file(deployment, out_dir + deployment_file)
#yaml.dump_to_file(service, out_dir + service_file)
yaml.dump_to_file(cnpg_cluster, out_dir + cnpg_cluster_file)
