package main

import (
	"flag"
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	rootDir := flag.String("root", "", "The root directory of this project")
	if *rootDir == "" {
		fmt.Println("❌ No root directory was specified as flag")
		return
	}
	name := ""
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "forgejo",
		GeneratorType: generatorType,
		ClusterUrl:    "forgejo-http.forgejo.svc.cluster.local",
		Port:          3000,
		Helm: &generator.Helm{
			Url:     "oci://code.forgejo.org/forgejo-helm/forgejo",
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "code.cluster",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 9 * * *",
			End:             "0 23 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{
			"postgres",
			"valkey",
		},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:      meta,
		OutputDir: filepath.Join(*rootDir, "/cluster/apps/forgejo/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createForgejoManifests(gm, *rootDir)
			if err != nil {
				fmt.Println("An error happened while generating Forgejo Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
