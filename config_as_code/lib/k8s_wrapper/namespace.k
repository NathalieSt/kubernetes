import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api

create_namespace = lambda name: str{
    core_api.Namespace{
    metadata: meta.ObjectMeta {
        name: name
        }
    }
}