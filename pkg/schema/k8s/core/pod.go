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
	Args         []string      `yaml:",omitempty"`
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

type SecurityContext struct {
	FsGroup    int `yaml:"fsGroup,omitempty"`
	RunAsUser  int `yaml:"runAsUser,omitempty"`
	RunAsGroup int `yaml:"runAsGroup,omitempty"`
}

type PodSpec struct {
	Containers      []Container     `yaml:",omitempty"`
	Volumes         []Volume        `yaml:",omitempty"`
	SecurityContext SecurityContext `yaml:"securityContext,omitempty"`
}

type Pod struct {
	shared.CommonK8sResourceWithSpec[PodSpec] `yaml:",omitempty,inline" validate:"required"`
}
