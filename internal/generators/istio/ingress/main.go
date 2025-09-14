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

	name := "istio-ingress"
	generatorType := generator.Istio
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "istio-ingress",
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Url:     "https://istio-release.storage.googleapis.com/charts",
			Chart:   "gateway",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/istio/ingress/"),
		CreateManifests:  createIngressManifests,
	})
}
