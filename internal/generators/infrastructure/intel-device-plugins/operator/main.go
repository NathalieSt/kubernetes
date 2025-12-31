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

	name := "intel-device-operator"
	generatorType := generator.Infrastructure
	var Vault = generator.GeneratorMeta{
		Name:          name,
		Namespace:     "inteldeviceplugins-system",
		GeneratorType: generatorType,
		Port:          8200,
		Helm: &generator.Helm{
			Chart:   "intel-device-plugins-operator",
			Url:     "https://intel.github.io/helm-charts/",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             Vault,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/intel-device-plugins/operator"),
		CreateManifests:  createVaultManifests,
	})
}
