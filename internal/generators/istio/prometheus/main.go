package main

import (
	"fmt"
	"kubernetes/internal/generators/istio"
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

	name := "prometheus"
	generatorType := generator.Istio
	Prometheus := generator.GeneratorMeta{
		Name:          name,
		Namespace:     istio.Namespace,
		GeneratorType: generatorType,
		Port:          9090,
		Docker: &generator.Docker{
			Registry: "prom/prometheus",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             Prometheus,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/istio/prometheus/"),
		CreateManifests:  createPrometheusManifests,
	})
}
