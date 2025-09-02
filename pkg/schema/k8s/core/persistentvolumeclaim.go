package core

import (
	"kubernetes/pkg/schema/k8s/meta"
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
	ApiVersion string `yaml:"apiVersion," validate:"required"`
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
