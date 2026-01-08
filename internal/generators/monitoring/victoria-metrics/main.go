package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/kustomization"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
		return
	}

	name := shared.VictoriaMetrics
	namespace := "victoria-metrics"
	generatorType := generator.Monitoring
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		ClusterUrl:    "vmsingle-victoria-metrics-vmks.victoria-metrics.svc.cluster.local",
		Port:          20001,
		Helm: &generator.Helm{
			Url:     "oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-k8s-stack",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/monitoring/victoria-metrics",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []string{
				shared.CSIDriverNFS,
			},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/victoria-metrics/"),
		CreateManifests:  createVictoriaMetricsManifests,
	})
}
