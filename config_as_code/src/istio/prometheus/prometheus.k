import ....lib.k8s_wrapper
import ...shared
import ..istio_shared
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml
import k8s.api.rbac.v1 as rbac_api
import outdent

out_dir = shared.cluster_path + "istio/prometheus/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
configmap_file = "configmap.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
rbac_file = "rbac.yaml"
pvc_file = "pvc.yaml"
# Configs
service_name = "prometheus"
label_app_name = "prometheus"
image_version = "v3.5.0"
k8s_service_name = "prometheus"
pvc_name = "prometheus-pvc"
configmap_name = "prometheus-configmap"
deployment_configmap_volume = "prometheus-config"
prometheus_storage_volume = "prometheus-storage"
istio_cert_volume = "istio-certs"
cluster_role_name = "istio-prometheus-cluster-role"
port = 9090
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    configmap_file
    deployment_file
    k8s_service_file
    rbac_file
    pvc_file
])
namespace = k8s_wrapper.create_namespace(service_name)

configmap = core_api.ConfigMap {
    apiVersion: "v1"
    kind: "ConfigMap"
    metadata: {
        name: configmap_name
        annotations: {
            "app.kubernetes.io/name": label_app_name
            "app.kubernetes.io/version": image_version
        }
    }
    data: {
        "prometheus.yml": outdent.outdent("""
global:
  scrape_interval: 30s
  evaluation_interval: 30s

scrape_configs:
  - job_name: 'istiod'
    kubernetes_sd_configs:
      - role: endpoints
        namespaces:
          names:
            - istio-system
    relabel_configs:
      - source_labels: [__meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
        action: keep
        regex: istiod;http-monitoring
    scheme: https
    tls_config:
      ca_file: /etc/prom-certs/root-cert.pem
      cert_file: /etc/prom-certs/cert-chain.pem
      key_file: /etc/prom-certs/key.pem
      insecure_skip_verify: true

  - job_name: 'envoy-stats'
    metrics_path: /stats/prometheus
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_container_port_name]
        action: keep
        regex: '.*-envoy-prom'
    scheme: https
    tls_config:
      ca_file: /etc/prom-certs/root-cert.pem
      cert_file: /etc/prom-certs/cert-chain.pem
      key_file: /etc/prom-certs/key.pem
      insecure_skip_verify: true
        """)
    }
}

cluster_role = rbac_api.ClusterRole {
    apiVersion = "rbac.authorization.k8s.io/v1"
    kind = "ClusterRole"
    metadata = {
        name = cluster_role_name
    }
    rules: [
        {
            apiGroups: [
                ""
            ]
            resources: [
                "nodes"
                "nodes/proxy"
                "services"
                "endpoints"
                "pods"
            ]
            verbs: [
                "get"
                "list"
                "watch"
            ]
        }
        {
            apiGroups: [
                "extensions"
            ]
            resources: [
                "ingresses"
            ]
            verbs = [
                "get"
                "list"
                "watch"
            ]
        }
        {
            nonResourceURLs: [
                "/metrics"
            ]
            verbs: [
                "get"
            ]
        }
    ]
}

cluster_role_binding: rbac_api.ClusterRoleBinding = {
    apiVersion = "rbac.authorization.k8s.io/v1"
    kind = "ClusterRoleBinding"
    metadata = {
        name = "istio-prometheus-cluster-role-binding"
    }
    roleRef = {
        apiGroup = "rbac.authorization.k8s.io"
        kind = "ClusterRole"
        name = cluster_role_name
    }
    subjects = [
        {
            kind = "ServiceAccount"
            name = "default"
            namespace = istio_shared.namespace
        }
    ]
}

rbac = [cluster_role, cluster_role_binding]

pvc: core_api.PersistentVolumeClaim = {
    metadata.name = pvc_name
    spec: {
        accessModes: ["ReadWriteMany"]
        storageClassName: shared.smb_storage_class_name
        resources.requests.storage: "100Gi"
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: "prometheus-deployment"
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
                    "sidecar.istio.io/inject": "true"
                }
                annotations: {
                    # do not intercept any inbound ports
                    "traffic.sidecar.istio.io/includeInboundPorts": ""
                    # do not intercept any outbound traffic
                    "traffic.sidecar.istio.io/includeOutboundIPRanges": ""
                    "proxy.istio.io/config": """
                        proxyMetadata:
                            OUTPUT_CERTS: /etc/istio-output-certs
                    """
                }
            }
            spec = core_api.PodSpec {
                containers: [
                    {
                        name: "prometheus"
                        image: "prom/prometheus:${image_version}"
                        args: [
                            "--storage.tsdb.retention.time=30d"
                            "--config.file=/etc/prometheus/prometheus.yml"
                            "--storage.tsdb.path=/prometheus/"
                        ]
                        ports: [
                            {
                                name: "prometheus"
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                mountPath: "/etc/prometheus"
                                name: deployment_configmap_volume
                            }
                            {
                                mountPath: "/prometheus/"
                                name: prometheus_storage_volume
                            }
                            {
                                mountPath: "/etc/prom-certs/"
                                name: istio_cert_volume
                            }
                        ]
                    }
                ]
                volumes: [
                    {
                        name: deployment_configmap_volume
                        configMap: {
                            name: configmap_name
                        }
                    }
                    {
                        name: prometheus_storage_volume
                        persistentVolumeClaim: {
                            claimName: pvc_name
                        }
                    }
                    {
                        name: istio_cert_volume
                        emptyDir: {
                            medium: "Memory"
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
                name: "prometheus"
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
yaml.dump_to_file(configmap, out_dir + configmap_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(service, out_dir + k8s_service_file)
yaml.dump_to_file(pvc, out_dir + pvc_file)
yaml.dump_all_to_file(rbac, out_dir + rbac_file)
