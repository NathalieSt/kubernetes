import ....lib.k8s_wrapper
import ...shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "apps/tandoor/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
staticfiles_pvc_file = "static-files-pvc.yaml"
mediafiles_pvc_file = "media-files-pvc.yaml"
# Configs
service_name = shared.tandoor.name
label_app_name = service_name
k8s_service_name = service_name
image_version = "2.0.1"
staticfiles_pvc_name = "static-files-" + service_name + "-pvc"
mediafiles_pvc_name = "media-files-" + service_name + "-pvc"
staticfiles_deployment_pvc_volume = "static-files-pvc-volume"
mediafiles_deployment_pvc_volume = "mediafiles-files-pvc-volume"
port = shared.tandoor.port
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    deployment_file
    k8s_service_file
    staticfiles_pvc_file
    mediafiles_pvc_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)

staticfiles_pvc: core_api.PersistentVolumeClaim = {
    apiVersion: "v1"
    kind: "PersistentVolumeClaim"
    metadata: {
        name: staticfiles_pvc_name
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_remote_storage_class_name
        resources.requests.storage: "5Gi"
    }
}

mediafiles_pvc: core_api.PersistentVolumeClaim = {
    apiVersion: "v1"
    kind: "PersistentVolumeClaim"
    metadata: {
        name: mediafiles_pvc_name
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.nfs_remote_storage_class_name
        resources.requests.storage: "5Gi"
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
                        name: "heimdall"
                        image: "linuxserver/heimdall:${image_version}"
                        ports: [
                            {
                                name: "http-${service_name}"
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                mountPath: "/opt/recipes/staticfiles"
                                name: staticfiles_deployment_pvc_volume
                            }
                            {
                                mountPath: "/opt/recipes/mediafiles"
                                name: mediafiles_deployment_pvc_volume
                            }
                        ]
                        env: [
                            {
                                name: "SECRET_KEY"
                                value: "4c9fMV%R1w!%*G0a&hfkb#2Tb6RYupZZ3s0z&xB^49yRQ!HwF7fQjN%uBV*X*8yN"
                            }
                            {
                                name: "POSTGRES_HOST"
                                value: shared.postgres.cluster_url
                            }
                            {
                                name: "POSTGRES_DB"
                                value: "tandoor-db"
                            }
                            {
                                name: "POSTGRES_PORT"
                                value: "${shared.postgres.port}"
                            }
                            {
                                name: "POSTGRES_USER"
                                value: "postgres"
                            }
                            {
                                name: "POSTGRES_PASSWORD"
                                valueFrom: {
                                    secretKeyRef: {
                                        name: shared.postgres_creds_secret_name
                                        key: "password"
                                    }
                                }
                            }
                        ]
                    }
                ]
                volumes: [
                    {
                        name: staticfiles_deployment_pvc_volume
                        persistentVolumeClaim: {
                            claimName: staticfiles_pvc_name
                        }
                    }
                    {
                        name: mediafiles_deployment_pvc_volume
                        persistentVolumeClaim: {
                            claimName: mediafiles_pvc_name
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
                name: "http-${service_name}"
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
yaml.dump_to_file(staticfiles_pvc, out_dir + staticfiles_pvc_file)
yaml.dump_to_file(mediafiles_pvc, out_dir + mediafiles_pvc_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(service, out_dir + k8s_service_file)
