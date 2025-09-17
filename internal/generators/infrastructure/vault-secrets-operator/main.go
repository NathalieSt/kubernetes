package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
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
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 8 * * *",
			End:             "0 22 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{
			"vault",
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/vault-secrets-operator/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createVaultSecretsOperatorManifests(flags.RootDir, gm)
			if err != nil {
				fmt.Println("An error happened while generating Vault-Secrets-Operator manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
