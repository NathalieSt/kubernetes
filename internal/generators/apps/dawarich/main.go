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

	name := ""
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "dawarich",
		GeneratorType: generatorType,
		ClusterUrl:    "dawarich.dawarich.svc.cluster.local",
		Port:          3000,
		Docker: &generator.Docker{
			Registry: "freikin/dawarich",
			Version:  utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "dawarich.cluster",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 8 * * *",
			End:             "0 22 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{
			"redis",
			"postgres",
		},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:      meta,
		OutputDir: filepath.Join(*rootDir, "/cluster/apps/dawarich/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createDawarichManifests(gm, *rootDir)
			if err != nil {
				fmt.Println("An error happened while generating Dawarich Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
