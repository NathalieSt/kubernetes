package main

import (
	"flag"
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {

	rootDir := flag.String("root", "", "The root directory of this project")
	if *rootDir == "" {
		fmt.Println("❌ No root directory was specified as flag")
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
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            meta,
		OutputDir:       filepath.Join(*rootDir, "/cluster/monitoring/victoria-metrics/"),
		CreateManifests: createVictoriaMetricsManifests,
	})
}
