package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"path/filepath"
)

func main() {
	fmt.Println("✅ Finding project root")
	rootDir, err := utils.FindRoot()
	if err != nil {
		fmt.Println("❌ An error occurred while finding the project root")
		fmt.Println("Error: " + err.Error())
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            Forgejo,
		OutputDir:       filepath.Join(rootDir, "/cluster/apps/forgejo/"),
		CreateManifests: createForgejoManifests,
	})
}
