package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
		return
	}

	name := "searxng"
	generatorType := generator.App
	var SearXNG = generator.GeneratorMeta{
		Name:          name,
		Namespace:     "searxng",
		GeneratorType: generatorType,
		ClusterUrl:    "searxng.searxng.svc.cluster.local",
		Port:          8080,
		Docker: &generator.Docker{
			Registry: "searxng/searxng",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "searxng",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 7 * * *",
			End:             "0 23 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{
			"valkey",
			"gluetun-proxy",
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             SearXNG,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/searxng/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createSearXNGManifests(gm, flags.RootDir)
			if err != nil {
				fmt.Println("An error happened while generating SearXNG Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
