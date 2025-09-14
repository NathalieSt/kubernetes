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

	name := "cnpg"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          "cnpg",
		Namespace:     "cnpg-system",
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Chart:   "cloudnative-pg",
			Url:     "https://cloudnative-pg.github.io/charts",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/cnpg/"),
		CreateManifests:  createCNPGManifests,
	})
}
