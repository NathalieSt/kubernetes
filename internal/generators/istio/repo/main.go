package main

import (
	"flag"
	"fmt"
	"kubernetes/internal/generators/istio"
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

	repo := generator.GeneratorMeta{
		Name:          "istio-repo",
		Namespace:     istio.Namespace,
		GeneratorType: generator.Istio,
		Helm: &generator.Helm{
			Url: "https://istio-release.storage.googleapis.com/charts",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            repo,
		OutputDir:       filepath.Join(*rootDir, "/cluster/istio/repo/"),
		CreateManifests: createRepoManifests,
	})
}
