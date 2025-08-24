import lib.flux
import ...shared
import ....lib.helpers
import yaml

out_dir = shared.cluster_path + "flux/"
# Filenames
apps_kustomizations_file = "apps.yaml"
infrastructure_kustomizations_file = "infrastructure.yaml"
istio_kustomizations_file = "istio.yaml"
monitoring_kustomizations_file = "monitoring.yaml"

# k8s configuration
app_kustomizations: [flux.FluxKustomization] = [helpers.create_cluster_kustomization(app_service.name, app_service.flux_kustomization.kustomization_path, app_service.namespace, app_service.flux_kustomization.dependencies) for app_service in shared.cluster_app_services]

infrastructure_kustomizations: [flux.FluxKustomization] = [helpers.create_cluster_kustomization(infrastructure_service.name, infrastructure_service.flux_kustomization.kustomization_path, infrastructure_service.namespace, infrastructure_service.flux_kustomization.dependencies) for infrastructure_service in shared.cluster_infrastructure_services]

istio_kustomizations: [flux.FluxKustomization] = [helpers.create_cluster_kustomization(istio_service.name, istio_service.flux_kustomization.kustomization_path, istio_service.namespace, istio_service.flux_kustomization.dependencies) for istio_service in shared.cluster_istio_services]

monitoring_kustomizations: [flux.FluxKustomization] = [helpers.create_cluster_kustomization(monitoring_service.name, monitoring_service.flux_kustomization.kustomization_path, monitoring_service.namespace, monitoring_service.flux_kustomization.dependencies) for monitoring_service in shared.cluster_monitoring_services]

# Write to files
yaml.dump_to_file(app_kustomizations, out_dir + apps_kustomizations_file)
yaml.dump_to_file(infrastructure_kustomizations, out_dir + infrastructure_kustomizations_file)
yaml.dump_to_file(istio_kustomizations, out_dir + istio_kustomizations_file)
yaml.dump_to_file(monitoring_kustomizations, out_dir + monitoring_kustomizations_file)
