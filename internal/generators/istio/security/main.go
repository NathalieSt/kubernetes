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

	security := generator.GeneratorMeta{
		Name:                "istio-security",
		Namespace:           istio.Namespace,
		GeneratorType:       generator.Istio,
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            security,
		OutputDir:       filepath.Join(*rootDir, "/cluster/istio/security/"),
		CreateManifests: createSecurityManifests,
	})
}
