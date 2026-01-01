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

	name := "piped-backend"
	generatorType := generator.Infrastructure
	var Piped = generator.GeneratorMeta{
		Name:          name,
		Namespace:     "piped",
		GeneratorType: generatorType,
		ClusterUrl:    "piped-backend.piped.svc.cluster.local",
		Port:          8200,
		Helm: &generator.Helm{
			Chart:   "piped",
			Url:     "https://helm.piped.video",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "piped-backend",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             Piped,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/piped/backend"),
		CreateManifests:  createPipedManifests,
	})
}
