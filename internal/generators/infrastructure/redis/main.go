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

	name := "redis"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "redis",
		GeneratorType: generatorType,
		ClusterUrl:    "redis.redis.svc.cluster.local",
		Port:          6379,
		Docker: &generator.Docker{
			Registry: "redis",
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
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/redis/"),
		CreateManifests:  createRedisManifests,
	})
}
