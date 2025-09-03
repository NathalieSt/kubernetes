package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type VolumeResourceRequirements struct {
	Limits   map[string]string `yaml:",omitempty"`
	Requests map[string]string `yaml:",omitempty"`
}

type PersistentVolumeClaimSpec struct {
	AccessModes      []string                   `yaml:"accessModes,omitempty"`
	Selector         meta.LabelSelector         `yaml:",omitempty"`
	Resources        VolumeResourceRequirements `yaml:",omitempty"`
	StorageClassName string                     `yaml:"storageClassName,omitempty"`
}

type PersistentVolumeClaim struct {
	shared.CommonK8sResourceWithSpec[PersistentVolumeClaimSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewPersistentVolumeClaim(meta meta.ObjectMeta, spec PersistentVolumeClaimSpec) PersistentVolumeClaim {
	return PersistentVolumeClaim{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[PersistentVolumeClaimSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "v1",
				Kind:       "PersistentVolumeClaim",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
