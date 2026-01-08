package core

import (
	"kubernetes/pkg/schema/shared"
)

type VolumeMount struct {
	MountPath string `yaml:"mountPath,omitempty"`
	SubPath   string `yaml:"subPath,omitempty"`
	Name      string `yaml:"name,omitempty"`
	Readonly  bool   `yaml:"readonly,omitempty"`
}

type Port struct {
	ContainerPort int64  `yaml:"containerPort,omitempty"`
	Name          string `yaml:"name,omitempty"`
}

type FieldRef struct {
	FieldPath string `yaml:"fieldPath,omitempty"`
}

type ConfigMapKeySelector struct {
	Key      string `yaml:"key,omitempty"`
	name     string `yaml:"name,omitempty"`
	optional bool   `yaml:"optional,omitempty"`
}

type ValueFrom struct {
	SecretKeyRef    SecretKeyRef         `yaml:"secretKeyRef,omitempty"`
	FieldRef        FieldRef             `yaml:"fieldRef,omitempty"`
	ConfigMapKeyRef ConfigMapKeySelector `yaml:"configMapKeyRef,omitempty"`
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

type DNSPolicy = string

const (
	ClusterFirst DNSPolicy = "ClusterFirst"
)

type PodToleration struct {
	Key string `yaml:"key,omitempty"`
	// FIXME: can be narrowed down to enum
	Operator string `yaml:"operator,omitempty"`
	Value    string `yaml:"value,omitempty"`
	// FIXME: can probably be narrowed down
	Effect string `yaml:"effect,omitempty"`
}

type MatchExpression struct {
	Key string `yaml:"key,omitempty"`
	// FIXME: can be narrowed down to enum
	Operator string   `yaml:"operator,omitempty"`
	Values   []string `yaml:"values,omitempty"`
}

type NodeSelectorTerm struct {
	MatchExpressions []MatchExpression `yaml:"matchExpressions,omitempty"`
}

type PodNodeRequiredDuringSchedulingIgnoredDuringExecution struct {
	NodeSelectorTerms []NodeSelectorTerm `yaml:"nodeSelectorTerms,omitempty"`
}

type PodNodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution PodNodeRequiredDuringSchedulingIgnoredDuringExecution `yaml:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

type PodAffinity struct {
	NodeAffinity PodNodeAffinity `yaml:"nodeAffinity,omitempty"`
}

type PodSpec struct {
	Affinity           PodAffinity        `yaml:"affinity,omitempty"`
	InitContainers     []Container        `yaml:"initContainers,omitempty"`
	ServiceAccountName string             `yaml:"serviceAccountName,omitempty"`
	DNSPolicy          DNSPolicy          `yaml:"dnsPolicy,omitempty"`
	Containers         []Container        `yaml:"containers,omitempty"`
	Volumes            []Volume           `yaml:"volumes,omitempty"`
	SecurityContext    PodSecurityContext `yaml:"securityContext,omitempty"`
	NodeSelector       map[string]string  `yaml:"nodeSelector,omitempty"`
	Tolerations        []PodToleration    `yaml:"tolerations,omitempty"`
}

type Pod struct {
	shared.CommonK8sResourceWithSpec[PodSpec] `yaml:",omitempty,inline" validate:"required"`
}
