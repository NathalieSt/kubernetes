package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createPrometheusManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	pvcName := fmt.Sprintf("%v-pvc", generatorMeta.Name)
	pvc := utils.ManifestConfig{
		Filename: "pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(
				meta.ObjectMeta{
					Name: pvcName,
				},
				core.PersistentVolumeClaimSpec{
					AccessModes: []string{"ReadWriteMany"},
					Resources: core.VolumeResourceRequirements{
						Requests: map[string]string{
							"storage": "60Gi",
						}},
					StorageClassName: generators.NFSRemoteClass,
				},
			),
		},
	}

	configMapName := "prometheus-configmap"
	configMap := utils.ManifestConfig{
		Filename: "configmap.yaml",
		Manifests: []any{
			core.NewConfigMap(
				meta.ObjectMeta{
					Name: configMapName,
				},
				map[string]string{
					"prometheus.yml": `
    global:
      scrape_interval: 30s
      evaluation_interval: 30s

    scrape_configs:
      - job_name: 'istiod'
        kubernetes_sd_configs:
          - role: endpoints
            namespaces:
              names:
                - istio-system
        relabel_configs:
          - source_labels: [__meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
            action: keep
            regex: istiod;http-monitoring
        scheme: http

      - job_name: 'envoy-stats'
        metrics_path: /stats/prometheus
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          - source_labels: [__meta_kubernetes_pod_container_port_name]
            action: keep
            regex: '.*-envoy-prom'
        scheme: http
        `,
				},
			),
		},
	}

	storageVolumeName := "mealie-pvc-volume"
	configMapVolumeName := "prometheus-configmap-volume"
	istioCertsVolumeName := "istio-certs"
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
								"sidecar.istio.io/inject":   "true",
							},
							Annotations: map[string]string{
								"traffic.sidecar.istio.io/includeInboundPorts":     "",
								"traffic.sidecar.istio.io/includeOutboundIPRanges": "",
								"proxy.istio.io/config": `proxyMetadata:
	OUTPUT_CERTS: /etc/istio-output-certs`,
								"sidecar.istio.io/userVolumeMount": "[{\"name\": \"istio-certs\", \"mountPath\": \"/etc/istio-output-certs\"}]",
							},
						},
						Spec: core.PodSpec{
							SecurityContext: core.SecurityContext{
								FsGroup:    65534,
								RunAsUser:  65534,
								RunAsGroup: 65534,
							},
							Containers: []core.Container{
								{
									Args: []string{
										"--storage.tsdb.retention.time=30d",
										"--config.file=/etc/prometheus/prometheus.yml",
										"--storage.tsdb.path=/prometheus/",
									},
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
											MountPath: "/etc/prometheus",
											Name:      configMapVolumeName,
										},
										{
											MountPath: "/prometheus/",
											Name:      storageVolumeName,
										},
										{
											MountPath: "/etc/prom-certs/",
											Name:      istioCertsVolumeName,
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: storageVolumeName,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: pvcName,
									},
								},
								{
									Name: configMapVolumeName,
									ConfigMap: core.ConfigMapVolumeSource{
										Name: configMapName,
									},
								},
								{
									Name: istioCertsVolumeName,
									EmptyDir: core.EmptyDirVolumeSource{
										Medium: core.Memory,
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
				},
				core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
					Ports: []core.ServicePort{
						{
							Name:       fmt.Sprintf("http-%v", generatorMeta.Name),
							Port:       9090,
							TargetPort: 9090,
						},
					},
				},
			),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			deployment.Filename,
			service.Filename,
			pvc.Filename,
			configMap.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, deployment, service, pvc, configMap})
}
