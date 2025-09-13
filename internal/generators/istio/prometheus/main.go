package main

import (
	"flag"
	"fmt"
	"kubernetes/internal/generators/istio"
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

	name := "prometheus"
	generatorType := generator.Istio
	Prometheus := generator.GeneratorMeta{
		Name:          name,
		Namespace:     istio.Namespace,
		GeneratorType: generatorType,
		Port:          9090,
		Docker: &generator.Docker{
			Registry: "prom/prometheus",
			Version:  utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            Prometheus,
		OutputDir:       filepath.Join(*rootDir, "/cluster/istio/prometheus/"),
		CreateManifests: createPrometheusManifests,
	})
}
