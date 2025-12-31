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

	name := "adguard-home"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "adguard-home",
		GeneratorType: generatorType,
		ClusterUrl:    "adguard-home.adguard-home.svc.cluster.local",
		Docker: &generator.Docker{
			Registry: "adguard/adguardhome",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/adguard-home/"),
		CreateManifests:  createBookloreManifests,
	})
}
