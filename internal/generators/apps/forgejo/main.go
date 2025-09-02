package main

import (
	"fmt"
	"kubernetes/internal/generators/apps"
	"kubernetes/internal/pkg/utils"
	"path/filepath"
)

func main() {
	fmt.Println("✅ Getting Meta for Forgejo")
	forgejoMeta := apps.Forgejo

	fmt.Println("✅ Finding project root")
	rootDir, err := utils.FindRoot()
	if err != nil {
		fmt.Println("❌ An error occurred while finding the project root")
		fmt.Println("Error: " + err.Error())
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            forgejoMeta,
		OutputDir:       filepath.Join(rootDir, "/cluster/apps/forgejo/"),
		CreateManifests: createForgejoManifests,
	})
}
