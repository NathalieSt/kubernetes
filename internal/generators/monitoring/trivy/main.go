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

	name := "trivy"
	generatorType := generator.Monitoring
	kiali := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "trivy",
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Chart:   "trivy-operator",
			Url:     "https://aquasecurity.github.io/helm-charts/",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             kiali,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/trivy/"),
		CreateManifests:  createKialiManifests,
	})
}
