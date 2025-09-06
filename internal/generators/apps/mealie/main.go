package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"path/filepath"
)

func main() {
	// FIXME: remove when done testing and optimizing? performance
	defer utils.Timer()()

	fmt.Println("✅ Finding project root")
	rootDir, err := utils.FindRoot()
	if err != nil {
		fmt.Println("❌ An error occurred while finding the project root")
		fmt.Println("Error: " + err.Error())
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            Mealie,
		OutputDir:       filepath.Join(rootDir, "/cluster/apps/mealie/"),
		CreateManifests: createMealieManifests,
	})
}
