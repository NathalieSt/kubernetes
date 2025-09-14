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

	name := "valkey"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "valkey",
		GeneratorType: generatorType,
		ClusterUrl:    "valkey.valkey.svc.cluster.local",
		Port:          6379,
		Docker: &generator.Docker{
			Registry: "valkey/valkey",
			Version:  "8-alpine3.22",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/valkey/"),
		CreateManifests:  createValkeyManifests,
	})
}
