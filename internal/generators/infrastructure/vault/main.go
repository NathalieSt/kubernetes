package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"path/filepath"
)

func main() {
	rootDir, err := utils.FindRoot()
	if err != nil {
		fmt.Println("❌ An error occurred while finding the project root")
		fmt.Println("Error: " + err.Error())
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            Vault,
		OutputDir:       filepath.Join(rootDir, "/cluster/infrastructure/vault/"),
		CreateManifests: createVaultManifests,
	})
}
