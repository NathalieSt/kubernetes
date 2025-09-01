package core

import (
	"kubernetes/pkg/schema/k8s/meta"
)

type PodTemplateSpec struct {
	Metadata meta.ObjectMeta
	Spec     PodSpec
}

type PodTemplate struct {
	Metadata meta.ObjectMeta
	Spec     PodTemplateSpec
}
