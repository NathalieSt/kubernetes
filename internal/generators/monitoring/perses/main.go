package main

import (
	"flag"
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {

	rootDir := flag.String("root", "", "The root directory of this project")
	if *rootDir == "" {
		fmt.Println("❌ No root directory was specified as flag")
		return
	}

	name := "perses"
	generatorType := generator.Monitoring

	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "perses",
		GeneratorType: generatorType,
		ClusterUrl:    "perses.perses.svc.cluster.local",
		Port:          8080,
		Docker: &generator.Docker{
			Registry: "persesdev/perses",
			Version:  utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "perses.cluster",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 8 * * *",
			End:             "0 22 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            meta,
		OutputDir:       filepath.Join(*rootDir, "/cluster/monitoring/perses/"),
		CreateManifests: createPersesManifests,
	})
}
