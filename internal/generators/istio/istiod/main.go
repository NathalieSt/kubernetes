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

	name := "istiod"
	generatorType := generator.Istio
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     istio.Namespace,
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Chart:   "istiod",
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            meta,
		OutputDir:       filepath.Join(*rootDir, "/cluster/istio/istiod/"),
		CreateManifests: createIstiodManifests,
	})
}
