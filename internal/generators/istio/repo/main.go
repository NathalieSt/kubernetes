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

	repo := generator.GeneratorMeta{
		Name:          "istio-repo",
		Namespace:     istio.Namespace,
		GeneratorType: generator.Istio,
		Helm: &generator.Helm{
			Url: "https://istio-release.storage.googleapis.com/charts",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             repo,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/istio/repo/"),
		CreateManifests:  createRepoManifests,
	})
}
