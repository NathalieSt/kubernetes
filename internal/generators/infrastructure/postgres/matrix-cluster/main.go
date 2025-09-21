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

	name := "matrix-pg"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "matrix-pg-cluster",
		GeneratorType: generatorType,
		ClusterUrl:    "matrix-pg-rw.matrix-pg-cluster.svc.cluster.local",
		Docker: &generator.Docker{
			Registry: "ghcr.io/cloudnative-pg/postgresql",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Port:                5432,
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/postgres/matrix"),
		CreateManifests:  createMatrixClusterManifests,
	})
}
