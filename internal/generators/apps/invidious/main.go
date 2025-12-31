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

	name := "invidious"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "invidious",
		GeneratorType: generatorType,
		ClusterUrl:    "invidious.invidious.svc.cluster.local",
		Port:          80,
		Caddy: &generator.Caddy{
			DNSName: "invidious",
		},
		Docker: &generator.Docker{
			Registry: "quay.io/invidious/invidious",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	relativeDir := "internal/generators/apps/invidious"
	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/invidious/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createInvidiousManifests(gm, flags.RootDir, relativeDir)
			if err != nil {
				fmt.Println("An error happened while generating Invidious manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
