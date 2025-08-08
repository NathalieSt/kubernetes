import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api

create_namespace = lambda name: str, inject_envoy: bool {
    envoy_label = "disabled"
    if inject_envoy == True:
        envoy_label = "enabled"

    core_api.Namespace {
        metadata: meta.ObjectMeta {
            name: name
            labels: {
                "istio-injection": envoy_label
            }
        }
    }
}
