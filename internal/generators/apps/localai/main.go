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

	name := "localai"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "localai",
		GeneratorType: generatorType,
		ClusterUrl:    "localai.localai.svc.cluster.local",
		Port:          8080,
		Helm: &generator.Helm{

			Url:     "https://go-skynet.github.io/helm-charts/",
			Chart:   "local-ai",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "ai",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 7 * * *",
			End:             "0 0 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/localai/"),
		CreateManifests:  createLocalAIManifests,
	})
}
