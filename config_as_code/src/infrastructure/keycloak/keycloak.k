import yaml
import ....lib.k8s_wrapper
import ...shared
import ...versions
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api

out_dir = shared.cluster_path + "infrastructure/keycloak/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
discovery_service_file = "discovery-service.yaml"
# Configs
service_name = shared.keycloak.name
label_app_name = service_name
k8s_service_name = service_name
port = shared.keycloak.port
image_version = versions.docker_images.keycloak
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    deployment_file
    k8s_service_file
    discovery_service_file
])
namespace = k8s_wrapper.create_namespace(service_name, True)

keycloak_env: [core_api.EnvVar] = [
    {name: "KC_BOOTSTRAP_ADMIN_USERNAME", value: "admin"}
    {name: "KC_BOOTSTRAP_ADMIN_PASSWORD", value: "admin"}
    {name: "KC_PROXY_HEADERS", value: "xforwarded"}
    {name: "KC_HTTP_ENABLED", value: "true"}
    {name: "KC_HOSTNAME_STRICT", value: "false"}
    {name: "KC_HEALTH_ENABLED", value: "true"}
    {name: "KC_CACHE", value: "ispn"}
    {name: "JAVA_OPTS_APPEND", value: "-Djgroups.bind.address=$(POD_IP)"}
    {name: "KC_DB_URL_DATABASE", value: "keycloak"}
    {name: "KC_DB_URL_HOST", value: "tcp://${shared.postgres.cluster_url}:${shared.postgres.port}"}
    {name: "KC_DB", value: "postgres"}
    {name: "KC_DB_PASSWORD", valueFrom: {
        secretKeyRef: {
            name: shared.postgres_creds_secret_name
            key: "password"
        }
    }}
    {name: "KC_DB_USERNAME", value: "postgres"}
]

deployment = {
    apiVersion: "apps/v1"
    kind: "StatefulSet"
    metadata: meta.ObjectMeta {
        name: service_name
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": versions.docker_images.keycloak
        }
    }
    spec: {
        replicas: 2
        selector: {
            matchLabels: {
                "app.kubernetes.io/name": label_app_name
            }
        }
        template: {
            metadata: meta.ObjectMeta {
                labels: {
                    "app.kubernetes.io/name": label_app_name
                    "app.kubernetes.io/version": versions.docker_images.keycloak
                }
            }
            spec = core_api.PodSpec {
                containers: [
                    {
                        name: "keycloak"
                        args: ["start"]
                        image: "quay.io/keycloak/keycloak:${image_version}"
                        ports: [
                            {
                                name: "keycloak"
                                containerPort: port
                            }
                        ]
                        env: keycloak_env
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
            "app.kubernetes.io/version": versions.docker_images.keycloak
        }
    }
    spec: {
        ports: [
            {
                name: "http-keycloak"
                port: port
                targetPort: port
            }
        ]
        selector: {
            "app.kubernetes.io/name": label_app_name
        }
    }
}

discovery_service: core_api.Service = {
    apiVersion: "v1"
    kind: "Service"
    metadata: {
        name: "keycloak-discovery"
        labels: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": versions.docker_images.keycloak
        }
    }
    spec: {
        clusterIP: "None"
        selector: {
            "app.kubernetes.io/name": label_app_name
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(k8s_service, out_dir + k8s_service_file)
yaml.dump_to_file(discovery_service, out_dir + discovery_service_file)
