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

	name := "postgres"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "postgres",
		GeneratorType: generatorType,
		ClusterUrl:    "postgres-rw.postgres.svc.cluster.local",
		Docker: &generator.Docker{
			Registry: "ghcr.io/cloudnative-pg/postgis",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Port:                5432,
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/postgres/"),
		CreateManifests:  createPostgresManifests,
	})
}
