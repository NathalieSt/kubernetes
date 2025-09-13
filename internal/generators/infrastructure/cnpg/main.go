package main

import (
	"flag"
	"fmt"
	"kubernetes/internal/pkg/utils"
	"path/filepath"
)

func main() {
	rootDir := flag.String("root", "", "The root directory of this project")
	if *rootDir == "" {
		fmt.Println("❌ No root directory was specified as flag")
		return
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            CNPG,
		OutputDir:       filepath.Join(rootDir, "/cluster/infrastructure/cnpg/"),
		CreateManifests: createCNPGManifests,
	})
}
