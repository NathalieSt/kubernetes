import yaml
import file
import shared.shared_variables

import .ingress
import .kustomization
import .namespace
import .release
import .repository

resources = {
    "namespace.yaml": yaml.encode(namespace.forgejo_namespace)
    "ingress.yaml": yaml.encode(ingress.forgejo_ingress)
    "kustomization.yaml": yaml.encode(kustomization.forgejo_kustomization)
    "release.yaml": yaml.encode(release.forgejo_release)
    "repository.yaml": yaml.encode(repository.forgejo_repository)
}

file.mkdir(shared_variables.cwd + "/apps/forgejo/")

_ = [file.write(shared_variables.cwd + "/apps/forgejo/" + k, v) for k, v in resources]
