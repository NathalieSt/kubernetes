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

	name := "synapse"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "synapse",
		GeneratorType: generatorType,
		ClusterUrl:    "synapse.synapse.svc.cluster.local",
		Port:          8008,
		Docker: &generator.Docker{
			Registry: "ghcr.io/element-hq/synapse",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "matrix.cluster",
		},
		DependsOnGenerators: []string{
			"postgres",
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/matrix/synapse"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createSynapseManifests(gm, flags.RootDir)
			if err != nil {
				fmt.Println("An error happened while generating Synapse Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
