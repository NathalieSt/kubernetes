package main

import (
	"flag"
	"fmt"
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

	name := "vault"
	generatorType := generator.Infrastructure
	var Vault = generator.GeneratorMeta{
		Name:          name,
		Namespace:     "vault",
		GeneratorType: generatorType,
		ClusterUrl:    "vault.vault.svc.cluster.local",
		Port:          8200,
		Helm: &generator.Helm{
			Chart:   "vault",
			Url:     "https://helm.releases.hashicorp.com",
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "vault.cluster",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            Vault,
		OutputDir:       filepath.Join(*rootDir, "/cluster/infrastructure/vault/"),
		CreateManifests: createVaultManifests,
	})
}
