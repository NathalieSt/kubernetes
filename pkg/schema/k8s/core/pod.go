package core

import (
	"kubernetes/pkg/schema/shared"
)

type VolumeMount struct {
	MountPath string `yaml:"mountPath,omitempty"`
	Name      string `yaml:"name,omitempty"`
}

type Port struct {
	ContainerPort int    `yaml:"containerPort,omitempty"`
	Name          string `yaml:"name,omitempty"`
}

type Env struct {
	Name      string    `yaml:"name,omitempty"`
	Value     string    `yaml:"value,omitempty"`
	ValueFrom ValueFrom `yaml:"valueFrom,omitempty"`
}
type Container struct {
	Args         []string      `yaml:"args,omitempty"`
	Command      []string      `yaml:"command,omitempty"`
	Image        string        `yaml:"image,omitempty"`
	Name         string        `yaml:"name,omitempty"`
	Ports        []Port        `yaml:"ports,omitempty"`
	VolumeMounts []VolumeMount `yaml:"volumeMounts,omitempty"`
	Env          []Env         `yaml:"env,omitempty"`
}

type SecretKeyRef struct {
	Key  string `yaml:"key,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type ValueFrom struct {
	SecretKeyRef SecretKeyRef `yaml:"secretKeyRef,omitempty"`
}

type Capabilities struct {
	Add []string `yaml:"add,omitempty"`
}

type SecurityContext struct {
	FsGroup      int          `yaml:"fsGroup,omitempty"`
	RunAsUser    int          `yaml:"runAsUser,omitempty"`
	RunAsGroup   int          `yaml:"runAsGroup,omitempty"`
	Capabilities Capabilities `yaml:"capabilities,omitempty"`
}

type PodSpec struct {
	Containers      []Container     `yaml:"containers,omitempty"`
	Volumes         []Volume        `yaml:"volumes,omitempty"`
	SecurityContext SecurityContext `yaml:"securityContext,omitempty"`
}

type Pod struct {
	shared.CommonK8sResourceWithSpec[PodSpec] `yaml:",omitempty,inline" validate:"required"`
}
