package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"sort"
	"strings"
)

func createCaddyManifests(rootDir string, generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	exposedGeneratorsMeta, err := utils.GetMetaForExposedGenerators(rootDir)
	if err != nil {
		fmt.Printf("An error happened while getting metadata for exposed services: \n %v", err)
	}

	manuallyDefinedMetas := []generator.GeneratorMeta{
		{
			ClusterUrl: "crowdsec-ui-agent.cloud.nathalie-stiefsohn.eu",
			Caddy: &generator.Caddy{
				DNSName: "crowdsec",
			},
			Port: 3000,
		},
	}

	exposedGeneratorsMeta = append(exposedGeneratorsMeta, manuallyDefinedMetas...)

	sort.Slice(exposedGeneratorsMeta, func(i, j int) bool {
		return exposedGeneratorsMeta[i].Name < exposedGeneratorsMeta[j].Name
	})

	configmapName := "caddy-configmap"
	configmap := utils.ManifestConfig{
		Filename: "configmap.yaml",
		Manifests: []any{
			getCaddyConfigMap(configmapName, exposedGeneratorsMeta),
		},
	}

	servicesDNSName := exposedGeneratorsMeta.GetDNSNames()

	pvcName := fmt.Sprintf("%v-pvc", generatorMeta.Name)
	pvc := utils.ManifestConfig{
		Filename: "pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: pvcName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "1Gi",
				}},
				StorageClassName: shared.NFSLocalClass,
			}),
		},
	}

	deployment := utils.ManifestConfig{
		Filename: "deployment.yaml",
		Manifests: []any{
			getDeployment(generatorMeta, configmapName, strings.Join(servicesDNSName, ","), pvcName),
		},
	}

	service := utils.ManifestConfig{
		Filename: "service.yaml",
		Manifests: []any{
			core.NewService(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				}, core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
					Ports: []core.ServicePort{
						{
							Name:       fmt.Sprintf("http-%v", generatorMeta.Name),
							Port:       generatorMeta.Port,
							TargetPort: generatorMeta.Port,
						},
					},
				},
			),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				namespace.Filename,
				deployment.Filename,
				service.Filename,
				configmap.Filename,
				pvc.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, configmap, deployment, service, pvc})
}
