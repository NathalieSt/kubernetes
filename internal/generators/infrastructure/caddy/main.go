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
	defer utils.Timer()()

	utils.RunGenerator(utils.GeneratorConfig{
		Meta: Caddy,
		// FIXME: maybe set the relative path to root in meta?
		OutputDir:       filepath.Join(rootDir, "/cluster/infrastructure/caddy/"),
		CreateManifests: createCaddyManifests,
	})
}
