package shared

import "kubernetes/pkg/schema/k8s/meta"

type CommonK8sResource struct {
	ApiVersion string          `yaml:"apiVersion," validate:"required"`
	Kind       string          `yaml:",omitempty" validate:"required"`
	Metadata   meta.ObjectMeta `yaml:",omitempty" validate:"required"`
}

type CommonK8sResourceWithSpec[T any] struct {
	CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
	Spec              T `yaml:",omitempty" validate:"required"`
}
