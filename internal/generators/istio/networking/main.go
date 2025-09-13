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

	name := "istio-networking"
	generatorType := generator.Istio
	networking := generator.GeneratorMeta{
		Name:          name,
		Namespace:     istio.Namespace,
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Url:     "oci://code.forgejo.org/forgejo-helm/forgejo",
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:      networking,
		OutputDir: filepath.Join(*rootDir, "/cluster/istio/networking/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createIstioNetworkingManifests(*rootDir, gm)
			if err != nil {
				fmt.Println("An error happened while generating Istio networking Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
