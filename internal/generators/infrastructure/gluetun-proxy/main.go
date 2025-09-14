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

	name := "gluetun-proxy"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "gluetun-proxy",
		GeneratorType: generatorType,
		ClusterUrl:    "gluetun-proxy.gluetun-proxy.svc.cluster.local",
		Port:          8888,
		Docker: &generator.Docker{
			Registry: "qmcgaw/gluetun",
			Version:  "v3.40",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/gluetun-proxy/"),
		CreateManifests:  createGluetunProxyManifests,
	})
}
