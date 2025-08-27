import ....lib.k8s_wrapper
import ....lib.cnpg
import ...shared
import yaml

out_dir = shared.cluster_path + "infrastructure/postgres/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
cnpg_cluster_file = "cnpg-cluster.yaml"
# Configs
service_name = shared.postgres.name
label_app_name = service_name
port = shared.postgres.port
# k8s configuration
kustomization = k8s_wrapper.create_kustomization(service_name, [
    namespace_file
    cnpg_cluster_file
])
namespace = k8s_wrapper.create_namespace(shared.postgres.namespace, True)

cnpg_cluster: cnpg.Cluster = {
    metadata: {
        name: service_name
    }
    spec: {
        instances: 3
        storage: {
            size: "15Gi"
            storageClass: shared.nfs_remote_storage_class_name
        }
        inheritedMetadata: {
            annotations: {
                "proxy.istio.io/config": '{"holdApplicationUntilProxyStarts": true}'
            }
        }
        enableSuperuserAccess: True
        superuserSecret: {
            name: shared.postgres_creds_secret_name
        }
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(cnpg_cluster, out_dir + cnpg_cluster_file)
