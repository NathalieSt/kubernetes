import k8s.api.core.v1 as core_api
import yaml

import ....lib.k8s_wrapper
import ....lib.helm_flux
import ...shared
import ..istio_shared

out_dir = shared.cluster_path + "istio/repo/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
repository_file = "repository.yaml"
# Configs
chart_url = "https://istio-release.storage.googleapis.com/charts"

namespace = k8s_wrapper.create_namespace(istio_shared.namespace)

kustomization = k8s_wrapper.create_kustomization(istio_shared.repo_name, [
    namespace_file
    repository_file
])

repo: helm_flux.HelmRepository = {
    metadata.name = istio_shared.repo_name
    kind: "HelmRepository"
    spec: {
        interval: "24h"
        url: chart_url
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
