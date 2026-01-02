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

	name := "transmission"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "jellyfin",
		GeneratorType: generatorType,
		ClusterUrl:    "transmission.jellyfin.svc.cluster.local",
		Port:          3000,
		Caddy: &generator.Caddy{
			DNSName: "transmission",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/mediaserver/transmission/"),
		CreateManifests:  createTransmissionManifests,
	})
}
