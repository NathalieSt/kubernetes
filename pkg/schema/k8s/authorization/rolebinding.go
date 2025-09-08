package authorization

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type RoleBinding struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
	RoleRef                  RoleRef   `yaml:"roleRef,omitempty"`
	Subjects                 []Subject `yaml:"subjects,omitempty"`
}

func NewRoleBinding(meta meta.ObjectMeta, roleRef RoleRef, subjects []Subject) RoleBinding {
	return RoleBinding{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "RoleBinding",
			Metadata:   meta,
		},
		RoleRef:  roleRef,
		Subjects: subjects,
	}
}
