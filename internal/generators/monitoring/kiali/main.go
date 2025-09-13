package main

import (
	"flag"
	"fmt"
	"kubernetes/internal/generators/istio"
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

	name := "kiali"
	generatorType := generator.Monitoring
	kiali := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "kiali-operator",
		GeneratorType: generatorType,
		ClusterUrl:    fmt.Sprintf("kiali.%v.svc.cluster.local", istio.Namespace),
		Port:          20001,
		Caddy: &generator.Caddy{
			DNSName: "kiali.cluster",
		},
		Helm: &generator.Helm{
			Chart:   "kiali-operator",
			Url:     "https://kiali.org/helm-charts",
			Version: utils.GetGeneratorVersionByType(*rootDir, name, generatorType),
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorConfig{
		Meta:            kiali,
		OutputDir:       filepath.Join(*rootDir, "/cluster/monitoring/kiali/"),
		CreateManifests: createKialiManifests,
	})
}
