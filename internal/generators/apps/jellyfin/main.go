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

	name := "jellyfin"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "jellyfin",
		GeneratorType: generatorType,
		ClusterUrl:    "jellyfin.jellyfin.svc.cluster.local",
		Port:          8096,
		Helm: &generator.Helm{
			Url:     "https://jellyfin.github.io/jellyfin-helm",
			Chart:   "jellyfin",
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName:                    "jellyfin.cluster",
			WebsocketSupportIsRequired: true,
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 9 * * *",
			End:             "0 23 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            meta,
		OutputDir:       filepath.Join(*rootDir, "/cluster/apps/jellyfin/"),
		CreateManifests: createJellyfinManifests,
	})
}
