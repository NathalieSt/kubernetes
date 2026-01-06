package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"strings"
)

func createNFSVolumesManifests(rootDir string, generatorMeta generator.GeneratorMeta) map[string][]byte {

	manifests := []utils.ManifestConfig{}

	discoveredGenerators, err := utils.GetDiscoveredGeneratorsMeta(rootDir)
	if err != nil {
		print("An error occurred while getting Metas of Discovered Generators")
		return nil
	}

	persistentVolumes := []utils.ManifestConfig{}
	fileNames := []string{}
	for _, generator := range discoveredGenerators {
		if generator.NFSVolumes != nil {
			volumesMap := generator.NFSVolumes
			for _, volume := range volumesMap {

				server := ""
				share := ""
				switch volume.StorageClass {
				case generators.NFSLocalClass:
					server = generators.NFSLocalServer
					share = generators.NFSLocalShare
				case generators.NFSRemoteClass:
					server = generators.NFSRemoteServer
					share = generators.NFSRemoteShare
				case generators.DebianStorageClass:
					server = generators.DebianServer
					share = generators.DebianShare
				}
				fileName := fmt.Sprintf("%v.yaml", volume.Name)
				fileNames = append(fileNames, fileName)
				persistentVolumes = append(persistentVolumes, utils.ManifestConfig{
					Filename: fileName,
					Manifests: []any{
						core.NewPersistentVolume(
							meta.ObjectMeta{
								Name: volume.Name,
							},
							core.PersistentVolumeSpec{
								AccessModes: []core.AccessModes{
									core.ReadWriteMany,
								},
								Capacity:                      volume.Capacity,
								PersistentVolumeReclaimPolicy: core.Retain,
								StorageClassName:              volume.StorageClass,
								CSIDriverNFS: core.CSIDriverNFS{
									Driver:       "nfs.csi.k8s.io",
									VolumeHandle: strings.Join([]string{server, share, fmt.Sprintf("/%v/%v/", generator.Name, volume.Name)}, ""),
									VolumeAttributes: core.CSIDriverNFSVolumeAttributes{
										Server: server,
										Share:  share,
									},
								},
							},
						),
					},
				},
				)
			}
		}
	}

	kustomization := utils.ManifestConfig{
		Filename:  "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, fileNames),
	}

	manifests = append(manifests, persistentVolumes...)
	manifests = append(manifests, kustomization)

	return utils.MarshalManifests(manifests)
}
