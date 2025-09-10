package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/storage"
)

func createCSIDriverNFSManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm, nil, nil)

	localStorageClass := utils.ManifestConfig{
		Filename: "local-storage-class-v2.yaml",
		Manifests: []any{
			storage.NewStorageClass(
				meta.ObjectMeta{
					Name: generators.NFSLocalClass,
				},
				storage.StorageClassData{
					Provisioner: "nfs.csi.k8s.io",
					Parameters: map[string]string{
						"server": "raspberry-pi-5-0",
						"share":  "/mnt/external_ssd",
					},
					ReclaimPolicy:        "Retain",
					VolumeBindingMode:    "Immediate",
					AllowVolumeExpansion: true,
					MountOptions: []string{
						"nfsvers=4.1",
					},
				}),
		},
	}

	remoteStorageClass := utils.ManifestConfig{
		Filename: "remote-storage-class.yaml",
		Manifests: []any{
			storage.NewStorageClass(
				meta.ObjectMeta{
					Name: generators.NFSRemoteClass,
				},
				storage.StorageClassData{
					Provisioner: "nfs.csi.k8s.io",
					Parameters: map[string]string{
						"server": "remote-fs.netbird.selfhosted",
						"share":  "/mnt/HC_Volume_103061115",
					},
					ReclaimPolicy:        "Retain",
					VolumeBindingMode:    "Immediate",
					AllowVolumeExpansion: true,
					MountOptions: []string{
						"nfsvers=4.1",
						"hard",
					},
				}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				namespace.Filename,
				repo.Filename,
				chart.Filename,
				release.Filename,
				localStorageClass.Filename,
				remoteStorageClass.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, namespace, repo, chart, release, release, localStorageClass, remoteStorageClass})
}
