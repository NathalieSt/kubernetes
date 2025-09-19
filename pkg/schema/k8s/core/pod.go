package core

import (
	"kubernetes/pkg/schema/shared"
)

type VolumeMount struct {
	MountPath string `yaml:"mountPath,omitempty"`
	Name      string `yaml:"name,omitempty"`
}

type Port struct {
	ContainerPort int64  `yaml:"containerPort,omitempty"`
	Name          string `yaml:"name,omitempty"`
}

type Env struct {
	Name      string    `yaml:"name,omitempty"`
	Value     string    `yaml:"value,omitempty"`
	ValueFrom ValueFrom `yaml:"valueFrom,omitempty"`
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

type ContainerSecurityContext struct {
	Capabilities Capabilities `yaml:"capabilities,omitempty"`
	Privileged   bool         `yaml:"privileged,omitempty"`
}

type PodSecurityContext struct {
	FsGroup    int `yaml:"fsGroup,omitempty"`
	RunAsUser  int `yaml:"runAsUser,omitempty"`
	RunAsGroup int `yaml:"runAsGroup,omitempty"`
}

type Resources struct {
	Requests map[string]string `yaml:"requests,omitempty"`
	Limits   map[string]string `yaml:"limits,omitempty"`
}

type Container struct {
	Args            []string                 `yaml:"args,omitempty"`
	Command         []string                 `yaml:"command,omitempty"`
	Image           string                   `yaml:"image,omitempty"`
	Name            string                   `yaml:"name,omitempty"`
	Ports           []Port                   `yaml:"ports,omitempty"`
	VolumeMounts    []VolumeMount            `yaml:"volumeMounts,omitempty"`
	Env             []Env                    `yaml:"env,omitempty"`
	Resources       Resources                `yaml:"resources,omitempty"`
	SecurityContext ContainerSecurityContext `yaml:"securityContext,omitempty"`
}

type PodSpec struct {
	Containers      []Container        `yaml:"containers,omitempty"`
	Volumes         []Volume           `yaml:"volumes,omitempty"`
	SecurityContext PodSecurityContext `yaml:"securityContext,omitempty"`
}

type Pod struct {
	shared.CommonK8sResourceWithSpec[PodSpec] `yaml:",omitempty,inline" validate:"required"`
}
