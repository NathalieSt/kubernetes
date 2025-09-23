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

	name := "whatsapp-bridge"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "whatsapp-bridge",
		GeneratorType: generatorType,
		ClusterUrl:    "whatsapp-bridge.whatsapp-bridge.svc.cluster.local",
		Port:          29318,
		Docker: &generator.Docker{
			Registry: "dock.mau.dev/mautrix/whatsapp",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		DependsOnGenerators: []string{
			"postgres",
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/matrix/whatsapp-bridge"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createWhatsappBridgeManifests(gm, flags.RootDir)
			if err != nil {
				fmt.Println("An error happened while generating Whatsapp Bridge Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
