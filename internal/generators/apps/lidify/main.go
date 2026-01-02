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

	name := "lidify"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "lidify",
		GeneratorType: generatorType,
		ClusterUrl:    "lidify.lidify.svc.cluster.local",
		Port:          3030,
		Docker: &generator.Docker{
			Registry: "chevron7locked/lidify",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "lidify",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 7 * * *",
			End:             "0 1 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/lidify/"),
		CreateManifests:  createLidifyManifests,
	})
}
