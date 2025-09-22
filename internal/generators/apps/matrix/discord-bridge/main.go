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
		Port:          29334,
		Docker: &generator.Docker{
			Registry: "dock.mau.dev/mautrix/discord",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{
			"postgres",
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/matrix/discord-bridge"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createDiscordBridgeManifests(gm, flags.RootDir)
			if err != nil {
				fmt.Println("An error happened while generating Discord Bridge Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
