import ....lib.k8s_wrapper
import ...shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "apps/heimdall/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
pvc_file = "pvc.yaml"
# Configs
service_name = shared.heimdall.name
label_app_name = service_name
k8s_service_name = service_name
image_version = "2.7.4"
pvc_name = service_name + "-pvc"
deployment_pvc_volume = "heimdall-config"
port = shared.heimdall.port
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    deployment_file
    k8s_service_file
    pvc_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)

pvc: core_api.PersistentVolumeClaim = {
    apiVersion: "v1"
    kind: "PersistentVolumeClaim"
    metadata: {
        name: pvc_name
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    spec: {
        accessModes: ["ReadWriteOnce"]
        storageClassName: "local-path"
        resources.requests.storage: "1Gi"
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
                        image: "linuxserver/heimdall:${image_version}"
                        env: [
                            {
                                name: "APP_URL"
                                value: "https://${shared.heimdall.gateway_config.dns_name}.netbird.selfhosted"
                            }
                        ]
                        ports: [
                            {
                                name: "http-${service_name}"
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                mountPath: "/config"
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
yaml.dump_to_file(pvc, out_dir + pvc_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(service, out_dir + k8s_service_file)
