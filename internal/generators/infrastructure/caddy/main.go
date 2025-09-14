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

	name := "caddy"
	generatorType := generator.Infrastructure
	caddy := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "caddy",
		GeneratorType: generatorType,
		ClusterUrl:    "caddy.caddy.svc.cluster.local",
		Port:          80,
		Docker: &generator.Docker{
			Registry: "caddy",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{
			"istio-networking",
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             caddy,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/caddy/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			return createCaddyManifests(flags.RootDir, gm)
		},
	})
}
