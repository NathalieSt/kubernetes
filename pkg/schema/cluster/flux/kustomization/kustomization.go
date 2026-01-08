package kustomization

/*
---
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: glance
  namespace: flux-system
spec:
  interval: "24h"
  targetNamespace: glance
  sourceRef:
    kind: GitRepository
    name: flux-system
  path: "./cluster/apps/glance"
  prune: true
  wait: true
  timeout: "10m"
  dependsOn: []
*/

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type SourceRefKind string

const (
	GitRepository    SourceRefKind = "GitRepository"
	OCIRepository    SourceRefKind = "OCIRepository"
	Bucket           SourceRefKind = "SourceRefKind"
	ExternalArtifact SourceRefKind = "SourceRefKind"
)

type SourceRef struct {
	Kind SourceRefKind `yaml:"kind,omitempty"`
	Name string        `yaml:"name,omitempty"`
}

type KustomizationSpec struct {
	Interval        string    `yaml:"interval,omitempty"`
	TargetNamespace string    `yaml:"targetNamespace,omitempty"`
	SourceRef       SourceRef `yaml:"sourceRef,omitempty"`
	Path            string    `yaml:"path,omitempty"`
	Prune           bool      `yaml:"prune,omitempty"`
	Wait            bool      `yaml:"wait,omitempty"`
	Timeout         string    `yaml:"timeout,omitempty"`
	DependsOn       []string  `yaml:"dependsOn,omitempty"`
}

type Kustomization struct {
	shared.CommonK8sResourceWithSpec[KustomizationSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewKustomization(meta meta.ObjectMeta, spec KustomizationSpec) Kustomization {
	return Kustomization{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[KustomizationSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "kustomize.toolkit.fluxcd.io/v1",
				Kind:       "Kustomization",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
