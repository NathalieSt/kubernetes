package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
		return
	}

	name := "victoria-metrics"
	generatorType := generator.Monitoring
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "victoria-metrics",
		GeneratorType: generatorType,
		ClusterUrl:    "vmsingle-victoria-metrics-vmks.victoria-metrics.svc.cluster.local",
		Port:          20001,
		Helm: &generator.Helm{
			Url:     "oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-k8s-stack",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/victoria-metrics/"),
		CreateManifests:  createVictoriaMetricsManifests,
	})
}
