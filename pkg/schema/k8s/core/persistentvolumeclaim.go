package core

import (
	"kubernetes/pkg/schema/k8s/meta"
)

type VolumeResourceRequirements struct {
	Limits   map[string]string
	Requests map[string]string
}

type PersistentVolumeClaimSpec struct {
	AccessModes      []string
	Selector         meta.LabelSelector
	Resources        VolumeResourceRequirements
	StorageClassName string
}

type PersistentVolumeClaim struct {
	ApiVersion string `validate:"required"`
	Kind       string `validate:"required"`
	Metadata   meta.ObjectMeta
	Spec       PersistentVolumeClaimSpec
}

func NewPersistentVolumeClaim(meta meta.ObjectMeta, spec PersistentVolumeClaimSpec) PersistentVolumeClaim {
	return PersistentVolumeClaim{
		ApiVersion: "v1",
		Kind:       "PersistentVolumeClaim",
		Metadata:   meta,
		Spec:       spec,
	}
}
