package authorization

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Rule struct {
	APIGroups       []string `yaml:"apiGroups,omitempty"`
	Resources       []string `yaml:"resources,omitempty"`
	Verbs           []string `yaml:"verbs,omitempty"`
	NonResourceURLs []string `yaml:"nonResourceURLs,omitempty"`
}

type ClusterRole struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
	Rules                    []Rule `yaml:"rules"`
}

func NewClusterRole(meta meta.ObjectMeta, rules []Rule) ClusterRole {
	return ClusterRole{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
			Metadata:   meta,
		},
		Rules: rules,
	}
}
