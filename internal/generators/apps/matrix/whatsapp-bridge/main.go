package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/kustomization"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
		return
	}

	name := shared.MatrixWhatsappBridge
	namespace := "whatsapp-bridge"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		ClusterUrl:    "whatsapp-bridge.whatsapp-bridge.svc.cluster.local",
		Port:          29318,
		Docker: &generator.Docker{
			Registry: "dock.mau.dev/mautrix/whatsapp",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/apps/matrix/whatsapp-bridge",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []string{
				shared.MatrixPostgres,
				shared.MatrixSynapse,
				shared.CSIDriverNFS,
			},
		},
	}

	relativeDir := "internal/generators/apps/matrix/whatsapp-bridge"
	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/matrix/whatsapp-bridge"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := createWhatsappBridgeManifests(gm, flags.RootDir, relativeDir)
			if err != nil {
				fmt.Println("An error happened while generating Whatsapp Bridge Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
