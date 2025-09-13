package main

import (
	"flag"
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	rootDir := flag.String("root", "", "The root directory of this project")
	if *rootDir == "" {
		fmt.Println("❌ No root directory was specified as flag")
		return
	}

	name := "caddy"
	generatorType := generator.Infrastructure
	caddy := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "caddy",
		GeneratorType: generatorType,
		ClusterUrl:    "caddy.caddy.svc.cluster.local",
		Port:          80,
		Docker: &generator.Docker{
			Registry: "caddy",
			Version:  utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		DependsOnGenerators: []string{
			"istio-networking",
		},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta: caddy,
		// FIXME: maybe set the relative path to root in meta?
		OutputDir: filepath.Join(*rootDir, "/cluster/infrastructure/caddy/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			return createCaddyManifests(*rootDir, gm)
		},
	})
}
