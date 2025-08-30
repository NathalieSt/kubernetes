import ....lib.k8s_wrapper
import ...shared
import ...versions
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import k8s.api.core.v1 as core_api
import yaml
import outdent
import ....lib.keda

out_dir = shared.cluster_path + "apps/glance/"
# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
configmap_file = "configmap.yaml"
deployment_file = "deployment.yaml"
k8s_service_file = "service.yaml"
scaled_object_file = "scaled-object.yaml"
# Configs
service_name = shared.glance.name
namespace_name = shared.glance.namespace
label_app_name = service_name
image_version = versions.docker_images.glance
k8s_service_name = service_name
configmap_name = "glance-configmap"
deployment_configmap_volume = "glance-config"
port = shared.glance.port
# k8s configuration
namespace = k8s_wrapper.create_namespace(service_name, True)

kustomization = k8s_wrapper.create_kustomization(service_name, [
    configmap_file
    deployment_file
    k8s_service_file
    namespace_file
    scaled_object_file
])

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
        "glance.yml": outdent.outdent("""
pages:
  - name: Home
    # Optionally, if you only have a single page you can hide the desktop navigation for a cleaner look
    # hide-desktop-navigation: true
    columns:
      - size: small
        widgets:
          - type: calendar
            first-day-of-week: monday

          - type: rss
            limit: 10
            collapse-after: 3
            cache: 12h
            feeds:
              - url: https://selfh.st/rss/
                title: selfh.st
                limit: 4
              - url: https://ciechanow.ski/atom.xml
              - url: https://www.joshwcomeau.com/rss.xml
                title: Josh Comeau
              - url: https://samwho.dev/rss.xml
              - url: https://ishadeed.com/feed.xml
                title: Ahmad Shadeed

          - type: twitch-channels
            channels:
              - theprimeagen
              - j_blow
              - piratesoftware
              - cohhcarnage
              - christitustech
              - EJ_SA

      - size: full
        widgets:
          - type: group
            widgets:
              - type: hacker-news
              - type: lobsters

          - type: videos
            channels:
              - UCXuqSBlHAE6Xw-yeJA0Tunw # Linus Tech Tips
              - UCR-DXc1voovS8nhAvccRZhg # Jeff Geerling
              - UCsBjURrPoezykLs9EqgamOA # Fireship
              - UCBJycsmduvYEL83R_U4JriQ # Marques Brownlee
              - UCHnyfMqiRRG1u-2MsSQLbXA # Veritasium

          - type: group
            widgets:
              - type: reddit
                subreddit: technology
                show-thumbnails: true
              - type: reddit
                subreddit: selfhosted
                show-thumbnails: true

      - size: small
        widgets:
          - type: weather
            location: London, United Kingdom
            units: metric # alternatively "imperial"
            hour-format: 12h # alternatively "24h"
            # Optionally hide the location from being displayed in the widget
            # hide-location: true

          - type: markets
            markets:
              - symbol: SPY
                name: S&P 500
              - symbol: BTC-USD
                name: Bitcoin
              - symbol: NVDA
                name: NVIDIA
              - symbol: AAPL
                name: Apple
              - symbol: MSFT
                name: Microsoft

          - type: releases
            cache: 1d
            # Without authentication the Github API allows for up to 60 requests per hour. You can create a
            # read-only token from your Github account settings and use it here to increase the limit.
            # token: ...
            repositories:
              - glanceapp/glance
              - go-gitea/gitea
              - immich-app/immich
              - syncthing/syncthing
        """)
    }
}

deployment = {
    apiVersion: "apps/v1"
    kind: "Deployment"
    metadata: meta.ObjectMeta {
        name: service_name
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
                }
            }
            spec = core_api.PodSpec {
                containers: [
                    {
                        name: service_name
                        image: "glanceapp/glance:${image_version}"
                        ports: [
                            {
                                name: service_name
                                containerPort: port
                            }
                        ]
                        volumeMounts: [
                            {
                                mountPath: "/app/config"
                                name: deployment_configmap_volume
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
                name: "http-${service_name}"
                port: port
                targetPort: port
            }
        ]
        selector: {
            "app.kubernetes.io/name": label_app_name
        }
    }
}

scaled_object: keda.ScaledObject = {
    metadata: {
        name: "${service_name}-scaledobject"
        namespace: "${namespace_name}"
    }
    spec: {
        minReplicaCount: 0
        cooldownPeriod: 300
        scaleTargetRef: {
            name: "${service_name}"
        }
        triggers: [
            {
                type: "cron"
                metadata: shared.glance.keda_scaling
            }
        ]
    }
}

# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(configmap, out_dir + configmap_file)
yaml.dump_to_file(deployment, out_dir + deployment_file)
yaml.dump_to_file(service, out_dir + k8s_service_file)
yaml.dump_to_file(scaled_object, out_dir + scaled_object_file)
