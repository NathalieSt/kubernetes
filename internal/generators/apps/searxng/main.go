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

	name := "searxng"
	generatorTy
	var SearXNG = generator.GeneratorMeta{
		Name:          "",
		Namespace:     "searxng",
		GeneratorType: generator.App,
		ClusterUrl:    "searxng.searxng.svc.cluster.local",
		Port:          8080,
		Docker: &generator.Docker{
			Registry: "searxng/searxng",
			//FIXME: set to nil, later fetch in generator from version.json
			Version: "2025.8.3-2e62eb5",
		},
		Caddy: &generator.Caddy{
			DNSName: "searxng.cluster",
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

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:      SearXNG,
		OutputDir: filepath.Join(rootDir, "/cluster/apps/searxng/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createSearXNGManifests(gm, rootDir)
			if err != nil {
				fmt.Println("An error happened while generating SearXNG Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
