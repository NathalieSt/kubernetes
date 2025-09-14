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

	name := "istiod"
	generatorType := generator.Istio
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     istio.Namespace,
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Chart:   "istiod",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/istio/istiod/"),
		CreateManifests:  createIstiodManifests,
	})
}
