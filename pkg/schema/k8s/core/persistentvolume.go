package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type AccessModes = string

const (
	ReadWriteOnce    AccessModes = "ReadWriteOnce"
	ReadOnlyMany     AccessModes = "ReadOnlyMany"
	ReadWriteMany    AccessModes = "ReadWriteMany"
	ReadWriteOncePod AccessModes = "ReadWriteOncePod"
)

type PersistentVolumeReclaimPolicy = string

const (
	Delete  PersistentVolumeReclaimPolicy = "Delete"
	Recycle PersistentVolumeReclaimPolicy = "Recycle"
	Retain  PersistentVolumeReclaimPolicy = "Retain"
)

type ObjectReference struct {
	ApiVersion      string `yaml:"apiVersion,omitempty"`
	FieldPath       string `yaml:"fieldPath,omitempty"`
	Kind            string `yaml:"kind,omitempty"`
	Name            string `yaml:"name,omitempty"`
	Namespace       string `yaml:"namespace,omitempty"`
	ResourceVersion string `yaml:"resourceVersion,omitempty"`
	UID             string `yaml:"uid,omitempty"`
}

type CSIDriverNFSVolumeAttributes struct {
	Server string `yaml:"server,omitempty"`
	Share  string `yaml:"share,omitempty"`
}

type CSIDriverNFS struct {
	Driver           string                       `yaml:"driver,omitempty"`
	VolumeHandle     string                       `yaml:"volumeHandle,omitempty"`
	VolumeAttributes CSIDriverNFSVolumeAttributes `yaml:"volumeAttributes,omitempty"`
}

type VolumeCapacity struct {
	Storage string `yaml:"storage,omitempty"`
}

type PersistentVolumeSpec struct {
	AccessModes                   []AccessModes                 `yaml:"accessModes,omitempty"`
	Capacity                      VolumeCapacity                `yaml:"capacity,omitempty"`
	ClaimRef                      ObjectReference               `yaml:"claimRef,omitempty"`
	MountOptions                  []string                      `yaml:"mountOptions,omitempty"`
	PersistentVolumeReclaimPolicy PersistentVolumeReclaimPolicy `yaml:"persistentVolumeReclaimPolicy,omitempty"`
	StorageClassName              string                        `yaml:"storageClassName,omitempty"`
	VolumeAttributesClassName     string                        `yaml:"volumeAttributesClassName,omitempty"`
	VolumeMode                    string                        `yaml:"volumeMode,omitempty"`
	CSIDriverNFS                  CSIDriverNFS                  `yaml:"csi,omitempty"`
}

type PersistentVolume struct {
	shared.CommonK8sResourceWithSpec[PersistentVolumeSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewPersistentVolume(meta meta.ObjectMeta, spec PersistentVolumeSpec) PersistentVolume {
	return PersistentVolume{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[PersistentVolumeSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "v1",
				Kind:       "PersistentVolume",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
