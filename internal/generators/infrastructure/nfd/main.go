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

	name := "nfd"
	generatorType := generator.Infrastructure
	var Vault = generator.GeneratorMeta{
		Name:          name,
		Namespace:     "nfd",
		GeneratorType: generatorType,
		Port:          8200,
		Helm: &generator.Helm{
			Chart:   "node-feature-discovery",
			Url:     "https://kubernetes-sigs.github.io/node-feature-discovery/charts",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             Vault,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/nfd"),
		CreateManifests:  createVaultManifests,
	})
}
