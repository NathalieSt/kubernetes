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
		Port:          80,
		Caddy: &generator.Caddy{
			DNSName: "adguard-home",
		},
		Docker: &generator.Docker{
			Registry: "adguard/adguardhome",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	relativeDir := "internal/generators/apps/matrix/discord-bridge"
	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/adguard-home/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createInvidiousManifests(gm, flags.RootDir, relativeDir)
			if err != nil {
				fmt.Println("An error happened while generating Discord Bridge Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
