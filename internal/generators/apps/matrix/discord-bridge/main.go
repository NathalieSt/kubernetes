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

	name := "discord-bridge"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "discord-bridge",
		GeneratorType: generatorType,
		ClusterUrl:    "discord-bridge.discord-bridge.svc.cluster.local",
		Port:          8008,
		Docker: &generator.Docker{
			Registry: "ghcr.io/mealie-recipes/mealie",
			//FIXME: set to nil, later fetch in generator from version.json
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{
			"postgres",
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/mealie/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createMealieManifests(gm, flags.RootDir)
			if err != nil {
				fmt.Println("An error happened while generating Forgejo Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
