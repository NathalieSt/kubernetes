package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	// FIXME: remove when done testing and optimizing? performance
	//defer utils.Timer()()

	rootDir, err := utils.FindRoot()
	if err != nil {
		fmt.Println("❌ An error occurred while finding the project root")
		fmt.Println("Error: " + err.Error())
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:      Mealie,
		OutputDir: filepath.Join(rootDir, "/cluster/apps/mealie/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createMealieManifests(gm, rootDir)
			if err != nil {
				fmt.Println("An error happened while generating Forgejo Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
