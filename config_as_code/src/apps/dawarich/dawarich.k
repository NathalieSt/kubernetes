import ....lib.k8s_wrapper
import ...shared
import ...versions
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "apps/dawarich/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
public_pvc_file = "public-dawarich-pvc.yaml"
watched_pvc_file = "watched-dawarich-pvc.yaml"
storage_pvc_file = "storage-dawarich-pvc.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
# Configs
service_name = shared.dawarich.name
label_app_name = "dawarich"
k8s_service_name = "dawarich"
public_pvc_name = "public-dawarich-pvc"
watched_pvc_name = "watched-dawarich-pvc"
storage_pvc_name = "storage-dawarich-pvc"
port = 3000
public_pvc_volume_name = "public-dawarich-volume"
watched_pvc_volume_name = "watched-dawarich-volume"
storage_pvc_volume_name = "storage-dawarich-volume"
image_version = versions.docker_images.dawarich
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    public_pvc_file
    watched_pvc_file
    storage_pvc_file
    deployment_file
    k8s_service_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)

public_pvc: core_api.PersistentVolumeClaim = {
    metadata.name = public_pvc_name
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_remote_storage_class_name
        resources.requests.storage: "10Gi"
    }
}

watched_pvc: core_api.PersistentVolumeClaim = {
    metadata.name = watched_pvc_name
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_remote_storage_class_name
        resources.requests.storage: "10Gi"
    }
}

storage_pvc: core_api.PersistentVolumeClaim = {
    metadata.name = storage_pvc_name
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_remote_storage_class_name
        resources.requests.storage: "10Gi"
    }
}

dawarich_env: [core_api.EnvVar] = [
    {name: "RAILS_ENV", value: "development"}
    {name: "REDIS_URL", value: "redis://${shared.redis.cluster_url}:${shared.redis.port}"}
    {name: "DATABASE_HOST", value: shared.postgres.cluster_url}
    {name: "DATABASE_PORT", value: "${shared.postgres.port}"}
    {name: "DATABASE_USERNAME", value: "postgres"}
    {name: "DATABASE_PASSWORD", valueFrom: {
        secretKeyRef: {
            name: shared.postgres_creds_secret_name
            key: "password"
        }
    }}
    {name: "DATABASE_NAME", value: "dawarich_development"}
    {name: "MIN_MINUTES_SPENT_IN_CITY", value: "60"}
    {name: "APPLICATION_HOSTS", value: "dawarich.cluster.netbird.selfhosted"}
    {name: "TIME_ZONE", value: "Europe/Vienna"}
    {name: "APPLICATION_PROTOCOL", value: "http"}
    {name: "PROMETHEUS_EXPORTER_ENABLED", value: "false"}
    {name: "PROMETHEUS_EXPORTER_HOST", value: "0.0.0.0"}
    {name: "PROMETHEUS_EXPORTER_PORT", value: "9394"}
    {name: "SELF_HOSTED", value: "true"}
    {name: "STORE_GEODATA", value: "true"}
]

sidekiq_env: [core_api.EnvVar] = [
    {name: "RAILS_ENV", value: "development"}
    {name: "REDIS_URL", value: "redis://${shared.redis.cluster_url}:${shared.redis.port}"}
    {name: "DATABASE_HOST", value: shared.postgres.cluster_url}
    {name: "DATABASE_PORT", value: "${shared.postgres.port}"}
    {name: "DATABASE_USERNAME", value: "postgres"}
    {name: "DATABASE_PASSWORD", valueFrom: {
        secretKeyRef: {
            name: shared.postgres_creds_secret_name
            key: "password"
        }
    }}
    {name: "DATABASE_NAME", value: "dawarich_development"}
    {name: "APPLICATION_HOSTS", value: "dawarich.cluster.netbird.selfhosted"}
    {name: "BACKGROUND_PROCESSING_CONCURRENCY", value: "10"}
    {name: "APPLICATION_PROTOCOL", value: "http"}
    {name: "PROMETHEUS_EXPORTER_ENABLED", value: "false"}
    {name: "PROMETHEUS_EXPORTER_HOST", value: "dawarich_app"}
    {name: "PROMETHEUS_EXPORTER_PORT", value: "9394"}
    {name: "SELF_HOSTED", value: "true"}
    {name: "STORE_GEODATA", value: "true"}
]

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
                        name: "dawarich"
                        image: "freikin/dawarich:${image_version}"
                        command: ['web-entrypoint.sh']
                        args: ['bin/rails', 'server', '-p', '3000', '-b', '::']
                        ports: [
                            {
                                name: "dawarich"
                                containerPort: port
                            }
                        ]
                        env: dawarich_env
                        volumeMounts: [
                            {
                                mountPath: "/var/app/public"
                                name: public_pvc_volume_name
                            }
                            {
                                mountPath: "/var/app/tmp/imports/watched"
                                name: watched_pvc_volume_name
                            }
                            {
                                mountPath: "/var/app/storage"
                                name: storage_pvc_volume_name
                            }
                        ]
                    }
                    {
                        name: "dawarich-sidekiq"
                        image: "freikin/dawarich:${image_version}"
                        command: ['sidekiq-entrypoint.sh']
                        ports: [
                            {
                                name: "dawarich"
                                containerPort: port
                            }
                        ]
                        env: sidekiq_env
                        volumeMounts: [
                            # FIXME: if it works, maybe unify it with above volumeMounts
                            {
                                mountPath: "/var/app/public"
                                name: public_pvc_volume_name
                            }
                            {
                                mountPath: "/var/app/tmp/imports/watched"
                                name: watched_pvc_volume_name
                            }
                            {
                                mountPath: "/var/app/storage"
                                name: storage_pvc_volume_name
                            }
                        ]
                    }
                ]
                volumes: [
                    {
                        name: public_pvc_volume_name
                        persistentVolumeClaim: {
                            claimName: public_pvc_name
                        }
                    }
                    {
                        name: watched_pvc_volume_name
                        persistentVolumeClaim: {
                            claimName: watched_pvc_name
                        }
                    }
                    {
                        name: storage_pvc_volume_name
                        persistentVolumeClaim: {
                            claimName: storage_pvc_name
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
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        ports: [
            {
                name: "http"
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
yaml.dump_to_file(public_pvc, out_dir + public_pvc_file)
yaml.dump_to_file(watched_pvc, out_dir + watched_pvc_file)
yaml.dump_to_file(storage_pvc, out_dir + storage_pvc_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(k8s_service, out_dir + k8s_service_file)
