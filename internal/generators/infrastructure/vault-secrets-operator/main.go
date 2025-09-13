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

	name := "vault-secrets-operator"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "vault-secrets-operator",
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Chart:   "vault-secrets-operator",
			Url:     "https://helm.releases.hashicorp.com",
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		DependsOnGenerators: []string{
			"vault",
		},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:      meta,
		OutputDir: filepath.Join(*rootDir, "/cluster/infrastructure/vault-secrets-operator/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createVaultSecretsOperatorManifests(*rootDir, gm)
			if err != nil {
				fmt.Println("An error happened while generating Vault-Secrets-Operator manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
