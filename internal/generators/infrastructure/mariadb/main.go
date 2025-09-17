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

	name := "mariadb"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "mariadb",
		GeneratorType: generatorType,
		ClusterUrl:    "mariadb.mariadb.svc.cluster.local",
		Port:          3306,
		Docker: &generator.Docker{
			Registry: "lscr.io/linuxserver/mariadb",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 6 * * *",
			End:             "0 0 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/mariadb/"),
		CreateManifests:  createGluetunProxyManifests,
	})
}
