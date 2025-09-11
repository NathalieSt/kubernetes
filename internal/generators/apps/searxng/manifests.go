package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/istio"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"path"
)

func createSearXNGManifests(generatorMeta generator.GeneratorMeta, rootDir string) (map[string][]byte, error) {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	valkeyMeta, err := utils.GetGeneratorMeta(path.Join(rootDir, "internal/generators/infrastructure/valkey"))
	if err != nil {
		fmt.Println("An error happened while getting valkey meta")
		return nil, err
	}

	proxyMeta, err := utils.GetGeneratorMeta(path.Join(rootDir, "internal/generators/infrastructure/gluetun-proxy"))
	if err != nil {
		fmt.Println("An error happened while getting gluetun-proxy meta")
		return nil, err
	}

	configmapName := fmt.Sprintf("%v-configmap", generatorMeta.Name)
	configmap := utils.ManifestConfig{
		Filename: "configmap.yaml",
		Manifests: []any{
			core.NewConfigMap(meta.ObjectMeta{
				Name: configmapName,
				Annotations: map[string]string{
					"app.kubernetes.io/name": generatorMeta.Name,
				},
			}, map[string]string{
				"settings.yml": fmt.Sprintf(`
# SearXNG settings
use_default_settings: true

general:
	debug: false
	instance_name: "Nathalies SearXNG"

search:
	safe_search: 2
	autocomplete: 'duckduckgo'
	formats:
		- html

server:
	secret_key: "mBrieVOfZQzc7"
	limiter: true
	image_proxy: true

valkey:
	url: valkey://%v:%v/0

outgoing:
	request_timeout: 4.0       # default timeout in seconds, can be override by engine
	max_request_timeout: 10.0  # the maximum timeout in seconds
	useragent_suffix: ""       # information like an email address to the administrator
	pool_connections: 100      # Maximum number of allowable connections, or null
								# for no limits. The default is 100.
	pool_maxsize: 10           # Number of allowable keep-alive connections, or null
								# to always allow. The default is 10.
	enable_http2: true         # See https://www.python-httpx.org/http2/
	proxies:
		all://:
			- http://%v:%v
			`, valkeyMeta.ClusterUrl, valkeyMeta.Port, proxyMeta.ClusterUrl, proxyMeta.Port),
			},
			),
		},
	}

	volumeName := "configmap-volume"
	deployment := utils.ManifestConfig{
		Filename: "deployment.yaml",
		Manifests: []any{
			apps.NewDeployment(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				},
				apps.DeploymentSpec{
					Replicas: 1,
					Selector: meta.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name": generatorMeta.Name,
						},
					},
					Template: core.PodTemplateSpec{
						Metadata: meta.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":    generatorMeta.Name,
								"app.kubernetes.io/version": generatorMeta.Docker.Version,
							},
						},
						Spec: core.PodSpec{
							Containers: []core.Container{
								{
									Name:  generatorMeta.Name,
									Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Ports: []core.Port{
										{
											ContainerPort: generatorMeta.Port,
											Name:          generatorMeta.Name,
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/etc/searxng",
											Name:      volumeName,
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: volumeName,
									ConfigMap: core.ConfigMapVolumeSource{
										Name: configmapName,
									},
								},
							},
						},
					},
				},
			),
		},
	}

	service := utils.ManifestConfig{
		Filename: "service.yaml",
		Manifests: []any{
			core.NewService(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				}, core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
					Ports: []core.ServicePort{
						{
							Name:       fmt.Sprintf("http-%v", generatorMeta.Name),
							Port:       generatorMeta.Port,
							TargetPort: generatorMeta.Port,
						},
					},
				},
			),
		},
	}

	scaledObject := utils.ManifestConfig{
		Filename:  "scaled-object.yaml",
		Manifests: utils.GenerateCronScaler(fmt.Sprintf("%v-scaledobject", generatorMeta.Name), generatorMeta.Name, generatorMeta.KedaScaling),
	}

	proxyServiceEntry := utils.ManifestConfig{
		Filename: "proxy-service-entry.yaml",
		Manifests: []any{
			istio.NewServiceEntry(
				meta.ObjectMeta{
					Name: "searxng-proxy-service-entry",
				},
				istio.ServiceEntrySpec{
					Hosts: []string{proxyMeta.ClusterUrl},
					Ports: []istio.ServiceEntryPorts{
						{
							Number:   proxyMeta.Port,
							Name:     "tcp",
							Protocol: "TCP",
						},
					},
					Location:   "MESH_EXTERNAL",
					Resolution: "DNS",
				}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			deployment.Filename,
			configmap.Filename,
			service.Filename,
			proxyServiceEntry.Filename,
			scaledObject.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, deployment, configmap, service, proxyServiceEntry, scaledObject}), nil
}
