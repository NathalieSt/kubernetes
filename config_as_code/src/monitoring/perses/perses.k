import ....lib.k8s_wrapper
import ....lib.keda
import ...shared
import ...versions

import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml

out_dir = shared.cluster_path + "monitoring/perses/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
scaled_object_file = "scaled-object.yaml"
# Configs
service_name = shared.perses.name
label_app_name = "perses"
k8s_service_name = "perses"
deployment_name = "${service_name}-deployment"
image_version = versions.docker_images.perses
container_registry = "persesdev/perses"
port = 8080
namespace_name = shared.perses.namespace
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    deployment_file
    k8s_service_file
    scaled_object_file
])
namespace = k8s_wrapper.create_namespace(namespace_name, True)

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: deployment_name
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
                        image: "${container_registry}:${image_version}"
                        ports: [
                            {
                                name: service_name
                                containerPort: port
                            }
                        ]
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

scaled_object: keda.ScaledObject = {
    metadata: {
        name: "${service_name}-scaledobject"
        namespace: "${namespace_name}"
    }
    spec: {
        minReplicaCount: 0
        cooldownPeriod: 300
        scaleTargetRef: {
            name: deployment_name
        }
        triggers: [
            {
                type: "cron"
                metadata: shared.perses.keda_scaling
            }
        ]
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(service, out_dir + k8s_service_file)
yaml.dump_to_file(scaled_object, out_dir + scaled_object_file)
