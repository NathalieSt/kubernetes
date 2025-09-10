package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	rootDir, err := utils.FindRoot()
	if err != nil {
		fmt.Println("An error occurred while finding the project root")
		fmt.Println("Error: " + err.Error())
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:      Networking,
		OutputDir: filepath.Join(rootDir, "/cluster/istio/networking/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createIstioNetworkingManifests(gm)
			if err != nil {
				fmt.Println("An error happened while generating Istio networking Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
