package helm

/*
---
apiVersion: source.toolkit.fluxcd.io/v1
kind: HelmChart
metadata:
  name: podinfo
  namespace: default
spec:
  interval: 5m0s
  chart: podinfo
  reconcileStrategy: ChartVersion
  sourceRef:
    kind: HelmRepository
    name: podinfo
  version: '5.*'
*/

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type ChartReconcileStrategy string

const (
	ChartVersion ChartReconcileStrategy = "ChartVersion"
	Revision     ChartReconcileStrategy = "Revision"
)

type ChartSourceRefKind string

const (
	HelmRepository ChartSourceRefKind = "HelmRepository"
	GitRepository  ChartSourceRefKind = "GitRepository"
	Bucket         ChartSourceRefKind = "Bucket"
)

type ChartSourceRef struct {
	Kind ChartSourceRefKind `yaml:",omitempty" validate:"required"`
	Name string             `yaml:",omitempty" validate:"required"`
}

type ChartSpec struct {
	Interval          string                 `yaml:",omitempty"`
	Chart             string                 `yaml:",omitempty" validate:"required"`
	ReconcileStrategy ChartReconcileStrategy `yaml:"reconcileStrategy,omitempty"`
	SourceRef         ChartSourceRef         `yaml:"sourceRef,omitempty" validate:"required"`
	Version           string                 `yaml:",omitempty" validate:"required"`
}

type Chart struct {
	shared.CommonK8sResourceWithSpec[ChartSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewChart(meta meta.ObjectMeta, spec ChartSpec) Chart {
	return Chart{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[ChartSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "source.toolkit.fluxcd.io/v1",
				Kind:       "HelmChart",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
