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
							},
						},
						"vmagent": map[string]any{
							"spec": map[string]any{
								"inlineScrapeConfig": `
# cilium-agent — port 9962, one pod per node
- job_name: cilium-agent
  kubernetes_sd_configs:
    - role: pod
      namespaces:
        names:
          - kube-system
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_k8s_app]
      regex: cilium
      action: keep
    - source_labels: [__address__]
      regex: '(.+):\d+'
      replacement: "${1}:9962"
      target_label: __address__
    - source_labels: [__meta_kubernetes_pod_node_name]
      target_label: node
    - source_labels: [__meta_kubernetes_pod_name]
      target_label: pod
    - source_labels: [__meta_kubernetes_namespace]
      target_label: namespace
    - target_label: job
      replacement: cilium-agent

# cilium-operator — port 9963
- job_name: cilium-operator
  kubernetes_sd_configs:
    - role: pod
      namespaces:
        names:
          - kube-system
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_name]
      regex: cilium-operator
      action: keep
    - source_labels: [__address__]
      regex: '(.+):\d+'
      replacement: "${1}:9963"
      target_label: __address__
    - source_labels: [__meta_kubernetes_pod_node_name]
      target_label: node
    - source_labels: [__meta_kubernetes_pod_name]
      target_label: pod
    - source_labels: [__meta_kubernetes_namespace]
      target_label: namespace
    - target_label: job
      replacement: cilium-operator
# hubble — port 9965
- job_name: hubble
  kubernetes_sd_configs:
    - role: pod
      namespaces:
        names:
          - kube-system
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_k8s_app]
      regex: cilium
      action: keep
    - source_labels: [__address__]
      regex: '(.+):\d+'
      replacement: "${1}:9965"
      target_label: __address__
    - source_labels: [__meta_kubernetes_pod_node_name]
      target_label: node
    - source_labels: [__meta_kubernetes_pod_name]
      target_label: pod
    - target_label: job
      replacement: hubble
`,
							},
						},
					},
				},
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
