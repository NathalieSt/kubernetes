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

	name := "glance"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "glance",
		GeneratorType: generatorType,
		ClusterUrl:    "glance.glance.svc.cluster.local",
		Port:          8080,
		Docker: &generator.Docker{
			Registry: "glanceapp/glance",
			//FIXME: set to nil, later fetch in generator from version.json
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "glance.cluster",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 7 * * *",
			End:             "0 23 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            meta,
		OutputDir:       filepath.Join(*rootDir, "/cluster/apps/glance/"),
		CreateManifests: createGlanceManifests,
	})
}
