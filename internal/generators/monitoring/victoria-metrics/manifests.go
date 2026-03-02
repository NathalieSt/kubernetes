package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/cluster/flux/oci"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createVictoriaMetricsManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repoName := fmt.Sprintf("%v-repo", generatorMeta.Name)
	repo := utils.ManifestConfig{
		Filename: "repo.yaml",
		Manifests: []any{
			oci.NewRepo(
				meta.ObjectMeta{
					Name: repoName,
				},
				oci.RepoSpec{
					Url:      generatorMeta.Helm.Url,
					Interval: "24h",
					Ref: oci.RepoRef{
						Tag: generatorMeta.Helm.Version,
					},
				}),
		},
	}

	// Helper to build the NFD labeldrop relabelConfig — reused across all kubelet scrape jobs
	nfdLabelDrop := []map[string]any{
		{
			"action": "labeldrop",
			// Drops ALL of: feature_node_kubernetes_io_*, beta_kubernetes_io_*,
			// gpu_intel_com_*, intel_feature_node_kubernetes_io_*, node_kubernetes_io_*
			// These are scheduler labels from NFD — useless for monitoring.
			"regex": "^(feature_node_kubernetes_io|beta_kubernetes_io|gpu_intel_com|intel_feature_node_kubernetes_io|node_kubernetes_io_instance_type).*",
		},
	}

	// cadvisor keep-list + NFD drop combined
	cadvisorRelabelConfigs := append(
		[]map[string]any{
			{
				"action":       "keep",
				"sourceLabels": []string{"__name__"},
				"regex": "container_cpu_usage_seconds_total|" +
					"container_cpu_cfs_periods_total|" +
					"container_cpu_cfs_throttled_periods_total|" +
					"container_cpu_cfs_throttled_seconds_total|" +
					"container_memory_rss|" +
					"container_memory_working_set_bytes|" +
					"container_memory_cache|" +
					"container_memory_swap|" +
					"container_memory_usage_bytes|" +
					"container_oom_events_total|" +
					"container_fs_reads_bytes_total|" +
					"container_fs_writes_bytes_total|" +
					"container_network_receive_bytes_total|" +
					"container_network_transmit_bytes_total|" +
					"container_spec_cpu_quota|" +
					"container_spec_cpu_period|" +
					"container_spec_memory_limit_bytes|" +
					"machine_cpu_cores|" +
					"machine_memory_bytes|" +
					"cadvisor_version_info",
			},
		},
		nfdLabelDrop...,
	)

	release := utils.ManifestConfig{
		Filename: "release.yaml",
		Manifests: []any{
			helm.NewRelease(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
				},
				helm.ReleaseSpec{
					ReleaseName: generatorMeta.Name,
					Interval:    "24h",
					Timeout:     "10m",
					ChartRef: helm.ReleaseChartRef{
						Kind: helm.OCIRepository,
						Name: repoName,
					},
					Install: helm.ReleaseInstall{
						Remediation: helm.ReleaseInstallRemediation{
							Retries: 3,
						},
					},
					Values: map[string]any{
						"nameOverride":      "vmks",
						"grafana":           map[string]any{"enabled": false},
						"defaultDashboards": map[string]any{"enabled": false},
						"vmsingle": map[string]any{
							"spec": map[string]any{
								"storage": map[string]any{
									"storageClassName": shared.NFSRemoteClass,
								},
								"extraArgs": map[string]any{
									"maxLabelsPerTimeseries": "120",
								},
							},
						},
						"prometheus-node-exporter": map[string]any{
							"prometheus": map[string]any{
								"monitor": map[string]any{
									"relabelings": []map[string]any{
										{
											"sourceLabels": []string{"__meta_kubernetes_node_name"},
											"targetLabel":  "node",
										},
										{
											"action": "labeldrop",
											"regex":  "^(feature_node_kubernetes_io|beta_kubernetes_io|gpu_intel_com|intel_feature_node_kubernetes_io|node_kubernetes_io_instance_type).*",
										},
									},
								},
							},
						},
						"vmagent": map[string]any{
							"spec": map[string]any{
								"inlineScrapeConfig": `
- job_name: cilium-agent
  kubernetes_sd_configs:
    - role: pod
      namespaces:
        names: [kube-system]
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_k8s_app]
      regex: cilium
      action: keep
    - source_labels: [__meta_kubernetes_pod_container_port_number]
      regex: "9962"
      action: keep
    - source_labels: [__meta_kubernetes_pod_node_name]
      target_label: node
    - source_labels: [__meta_kubernetes_pod_name]
      target_label: pod
    - source_labels: [__meta_kubernetes_namespace]
      target_label: namespace
    - target_label: job
      replacement: cilium-agent

- job_name: cilium-operator
  kubernetes_sd_configs:
    - role: pod
      namespaces:
        names: [kube-system]
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_name]
      regex: cilium-operator
      action: keep
    - source_labels: [__meta_kubernetes_pod_container_port_number]
      regex: "9963"
      action: keep
    - source_labels: [__meta_kubernetes_pod_node_name]
      target_label: node
    - source_labels: [__meta_kubernetes_pod_name]
      target_label: pod
    - source_labels: [__meta_kubernetes_namespace]
      target_label: namespace
    - target_label: job
      replacement: cilium-operator
`},
						},

						"kubelet": map[string]any{
							"enabled": true,
							"vmScrapes": map[string]any{
								"cadvisor": map[string]any{
									"enabled": true,
									"spec": map[string]any{
										"path":                 "/metrics/cadvisor",
										"scheme":               "https",
										"honorLabels":          true,
										"bearerTokenFile":      "/var/run/secrets/kubernetes.io/serviceaccount/token",
										"tlsConfig":            map[string]any{"insecureSkipVerify": true},
										"metricRelabelConfigs": cadvisorRelabelConfigs,
									},
								},
								"resource": map[string]any{
									"enabled": true,
									"spec": map[string]any{
										"path":                 "/metrics/resource",
										"scheme":               "https",
										"honorLabels":          true,
										"bearerTokenFile":      "/var/run/secrets/kubernetes.io/serviceaccount/token",
										"tlsConfig":            map[string]any{"insecureSkipVerify": true},
										"metricRelabelConfigs": nfdLabelDrop,
									},
								},
								"kubelet": map[string]any{
									"enabled": true,
									"spec": map[string]any{
										"scheme":               "https",
										"honorLabels":          true,
										"bearerTokenFile":      "/var/run/secrets/kubernetes.io/serviceaccount/token",
										"tlsConfig":            map[string]any{"insecureSkipVerify": true},
										"metricRelabelConfigs": nfdLabelDrop,
									},
								},
								"probes": map[string]any{
									"enabled": true,
									"spec": map[string]any{
										"path":                 "/metrics/probes",
										"scheme":               "https",
										"honorLabels":          true,
										"bearerTokenFile":      "/var/run/secrets/kubernetes.io/serviceaccount/token",
										"tlsConfig":            map[string]any{"insecureSkipVerify": true},
										"metricRelabelConfigs": nfdLabelDrop,
									},
								},
							},
						},
					}},
			),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			repo.Filename,
			release.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, release})
}
