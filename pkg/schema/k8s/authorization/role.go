package authorization

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Role struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
	Rules                    []Rule `yaml:"rules,omitempty"`
}

func NewRole(meta meta.ObjectMeta, rules []Rule) Role {
	return Role{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
			Metadata:   meta,
		},
		Rules: rules,
	}
}
