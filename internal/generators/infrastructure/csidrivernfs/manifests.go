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
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm, map[string]any{
		"kubeletDir": "/var/lib/k0s/kubelet",
	}, nil)

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
						"server": generators.NFSLocalServer,
						"share":  generators.NFSLocalShare,
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

	localStorageClassNext := utils.ManifestConfig{
		Filename: "local-storage-class-next.yaml",
		Manifests: []any{
			storage.NewStorageClass(
				meta.ObjectMeta{
					Name: generators.NFSLocalClassNext,
				},
				storage.StorageClassData{
					Provisioner: "nfs.csi.k8s.io",
					Parameters: map[string]string{
						"server": generators.NFSLocalServerNext,
						"share":  generators.NFSLocalShareNext,
						"subDir": "${pvc.metadata.namespace}/${pvc.metadata.name}",
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

	debianStorageClass := utils.ManifestConfig{
		Filename: "debian-storage-class.yaml",
		Manifests: []any{
			storage.NewStorageClass(
				meta.ObjectMeta{
					Name: generators.DebianStorageClass,
				},
				storage.StorageClassData{
					Provisioner: "nfs.csi.k8s.io",
					Parameters: map[string]string{
						"server": generators.DebianServer,
						"share":  generators.DebianShare,
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
						"server": generators.NFSRemoteServer,
						"share":  generators.NFSRemoteShare,
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
				debianStorageClass.Filename,
				localStorageClassNext.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, namespace, repo, chart, release, release, localStorageClass, remoteStorageClass, debianStorageClass, localStorageClassNext})
}
