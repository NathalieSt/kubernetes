import yaml
import file
import shared.shared_variables

import .deployment
import .ingress
import .configmap
import .rbac
import .service
import .kustomization

resources = {
    "rbac.yaml": yaml.encode_all(rbac.homepage_rbac)
    "deployment.yaml": yaml.encode(deployment.homepage_deployment)
    "service.yaml": yaml.encode(service.homepage_service)
    "configmap.yaml": yaml.encode(configmap.homepage_configmap)
    "ingress.yaml": yaml.encode(ingress.homepage_ingress)
    "kustomization.yaml": yaml.encode(kustomization.homepage_kustomization)
}

file.mkdir(shared_variables.cwd + "/apps/homepage/")

_ = [file.write(shared_variables.cwd + "/apps/homepage/" + k, v) for k, v in resources]
