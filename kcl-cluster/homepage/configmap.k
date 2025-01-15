homepage_configmap = {
    apiVersion = "v1"
    kind = "ConfigMap"
    metadata = {
        name = "homepage"
        labels = {
            "app.kubernetes.io/name" = "homepage"
        }
    }
    data = {
        "kubernetes.yaml" = r"""mode: cluster
"""
        "settings.yaml" = ""
        "custom.css" = ""
        "custom.js" = ""
        "bookmarks.yaml" = r"""- Developer:
    - Github:
        - abbr: GH
          href: https://github.com/
"""
        "services.yaml" = r"""- My First Group:
    - My First Service:
        href: http://localhost/
        description: Homepage is awesome

- My Second Group:
    - My Second Service:
        href: http://localhost/
        description: Homepage is the best

- My Third Group:
    - My Third Service:
        href: http://localhost/
        description: Homepage is 😎
"""
        "widgets.yaml" = r"""- kubernetes:
    cluster:
      show: true
      cpu: true
      memory: true
      showLabel: true
      label: "cluster"
    nodes:
      show: true
      cpu: true
      memory: true
      showLabel: true
- resources:
    backend: resources
    expanded: true
    cpu: true
    memory: true
    network: default
- search:
    provider: duckduckgo
    target: _blank
"""
        "docker.yaml" = ""
    }
}
