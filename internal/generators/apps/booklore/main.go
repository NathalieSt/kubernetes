package main

import (
	"fmt"
	"kubernetes/internal/generators"
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

	name := "booklore"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "booklore",
		GeneratorType: generatorType,
		ClusterUrl:    "booklore.booklore.svc.cluster.local",
		Port:          6060,
		Docker: &generator.Docker{
			Registry: "booklore/booklore",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "booklore",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 7 * * *",
			End:             "0 23 * * *",
			DesiredReplicas: "1",
		},
		NFSVolumes: map[string]generator.GeneratorNFSVolume{
			"books": {
				Name:         "books-pv",
				StorageClass: generators.NFSLocalClass,
				Capacity:     "100Gi",
			},
			"data": {
				Name:         "data-pv",
				StorageClass: generators.NFSRemoteClass,
				Capacity:     "1Gi",
			},
			"bookdrop": {
				Name:         "bookdrop-pv",
				StorageClass: generators.NFSLocalClass,
				Capacity:     "10Gi",
			},
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/booklore/"),
		CreateManifests:  createBookloreManifests,
	})
}
