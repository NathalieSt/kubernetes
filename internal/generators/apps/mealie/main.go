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

	name := "mealie"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "mealie",
		GeneratorType: generatorType,
		ClusterUrl:    "mealie.mealie.svc.cluster.local",
		Port:          9000,
		Docker: &generator.Docker{
			Registry: "ghcr.io/mealie-recipes/mealie",
			//FIXME: set to nil, later fetch in generator from version.json
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "mealie.cluster",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 9 * * *",
			End:             "0 21 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{
			"postgres",
		},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:      meta,
		OutputDir: filepath.Join(*rootDir, "/cluster/apps/mealie/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createMealieManifests(gm, *rootDir)
			if err != nil {
				fmt.Println("An error happened while generating Forgejo Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
