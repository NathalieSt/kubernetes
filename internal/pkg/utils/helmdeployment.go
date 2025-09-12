package utils

import (
	"fmt"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func GetGenericRepoManifest(repoName string, generatorHelm *generator.Helm) ManifestConfig {

	return ManifestConfig{
		Filename: "repo.yaml",
		Manifests: []any{
			helm.NewRepo(meta.ObjectMeta{
				Name: repoName,
			},
				helm.RepoSpec{
					RepoType: helm.Default,
					Url:      generatorHelm.Url,
					Interval: "24h",
				}),
		},
	}
}

func GetGenericChartManifest(chartName string, generatorHelm *generator.Helm, repoName string) ManifestConfig {
	return ManifestConfig{
		Filename: "chart.yaml",
		Manifests: []any{
			helm.NewChart(meta.ObjectMeta{
				Name: chartName,
			}, helm.ChartSpec{
				Interval:          "24h",
				Chart:             generatorHelm.Chart,
				ReconcileStrategy: helm.ChartVersion,
				SourceRef: helm.ChartSourceRef{
					Kind: helm.HelmRepository,
					Name: repoName,
				},
				Version: generatorHelm.Version,
			}),
		},
	}
}

func GetGenericReleaseManifest(generatorName string, chartName string, values map[string]any, valuesFrom []helm.ReleaseValuesFrom) ManifestConfig {
	return ManifestConfig{
		Filename: "release.yaml",
		Manifests: []any{
			helm.NewRelease(meta.ObjectMeta{
				Name: generatorName,
			},
				helm.ReleaseSpec{
					ReleaseName: generatorName,
					Interval:    "24h",
					Timeout:     "10m",
					ChartRef: helm.ReleaseChartRef{
						Kind: helm.HelmChart,
						Name: chartName,
					},
					Install: helm.ReleaseInstall{
						Remediation: helm.ReleaseInstallRemediation{
							Retries: 3,
						},
					},
					ValuesFrom: valuesFrom,
					Values:     values,
				}),
		},
	}
}

func GetGenericHelmDeploymentManifests(
	generatorName string,
	generatorHelm *generator.Helm,
	values map[string]any,
	valuesFrom []helm.ReleaseValuesFrom,
) (ManifestConfig, ManifestConfig, ManifestConfig) {
	repoName := fmt.Sprintf("%v-repo", generatorName)
	repo := GetGenericRepoManifest(repoName, generatorHelm)

	chartName := fmt.Sprintf("%v-chart", generatorName)
	chart := GetGenericChartManifest(chartName, generatorHelm, repoName)

	release := GetGenericReleaseManifest(generatorName, chartName, values, valuesFrom)

	return repo, chart, release
}
