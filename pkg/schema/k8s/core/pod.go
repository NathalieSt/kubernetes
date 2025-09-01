package core

import (
	"kubernetes/pkg/schema/k8s/meta"
)

type VolumeMount struct {
	MountPath string
	Name      string
}

type Port struct {
	ContainerPort int
	Name          string
}

type Container struct {
	Image        string
	Name         string
	Ports        []Port
	VolumeMounts []VolumeMount
}

type PodSpec struct {
	Containers Container
	Volumes    []Volume
}

type Pod struct {
	ApiVersion string
	Kind       string
	Metadata   meta.ObjectMeta
	Spec       PodSpec
}
