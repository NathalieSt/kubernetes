package authorization

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Metadata struct {
	Name string `yaml:"name,omitempty"`
}
type RoleRef struct {
	APIGroup string `yaml:"apiGroup,omitempty"`
	Kind     string `yaml:"kind,omitempty"`
	Name     string `yaml:"name,omitempty"`
}
type Subject struct {
	Kind      string `yaml:"kind,omitempty"`
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}

type ClusterRoleBinding struct {
	shared.CommonK8sResource `yaml:",omitempty,inline"`
	RoleRef                  RoleRef   `yaml:"roleRef,omitempty"`
	Subjects                 []Subject `yaml:"subjects,omitempty"`
}

func NewClusterRoleBinding(meta meta.ObjectMeta, roleRef RoleRef, subjects []Subject) ClusterRoleBinding {
	return ClusterRoleBinding{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
			Metadata:   meta,
		},
		RoleRef:  roleRef,
		Subjects: subjects,
	}
}
