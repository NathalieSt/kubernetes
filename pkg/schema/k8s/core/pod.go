package core

import (
	"kubernetes/pkg/schema/shared"
)

type VolumeMount struct {
	MountPath string `yaml:"mountPath,omitempty"`
	Name      string `yaml:",omitempty"`
}

type Port struct {
	ContainerPort int    `yaml:"containerPort,omitempty"`
	Name          string `yaml:",omitempty"`
}

type Env struct {
	Name      string    `yaml:",omitempty"`
	Value     string    `yaml:",omitempty"`
	ValueFrom ValueFrom `yaml:"valueFrom,omitempty"`
}
type Container struct {
	Image        string        `yaml:",omitempty"`
	Name         string        `yaml:",omitempty"`
	Ports        []Port        `yaml:",omitempty"`
	VolumeMounts []VolumeMount `yaml:"volumeMounts,omitempty"`
	Env          []Env         `yaml:",omitempty"`
}

type SecretKeyRef struct {
	Key  string `yaml:",omitempty"`
	Name string `yaml:",omitempty"`
}

type ValueFrom struct {
	SecretKeyRef SecretKeyRef `yaml:"secretKeyRef,omitempty"`
}

type PodSpec struct {
	Containers []Container `yaml:",omitempty"`
	Volumes    []Volume    `yaml:",omitempty"`
}

type Pod struct {
	shared.CommonK8sResourceWithSpec[PodSpec] `yaml:",omitempty,inline" validate:"required"`
}
