package helm

import "kubernetes/pkg/schema/k8s/meta"

type ReleaseChartRefKind string

const (
	OCIRepository ReleaseChartRefKind = "OCIRepository"
	HelmChart     ReleaseChartRefKind = "HelmChart"
)

type ReleaseChartRef struct {
	Kind      ReleaseChartRefKind `yaml:",omitempty" validate:"required"`
	Name      string              `yaml:",omitempty" validate:"required"`
	Namespace string              `yaml:",omitempty"`
}

type ReleaseInstallRemediation struct {
	Retries int `yaml:",omitempty"`
}

type ReleaseInstall struct {
	Remediation ReleaseInstallRemediation `yaml:",omitempty"`
}

type ReleaseValuesFromKind string

const (
	ConfigMap ReleaseValuesFromKind = "ConfigMap"
	Secret    ReleaseValuesFromKind = "Secret"
)

type ReleaseValuesFrom struct {
	Kind       ReleaseValuesFromKind `yaml:",omitempty" validate:"required"`
	Name       string                `yaml:",omitempty" validate:"required"`
	ValuesKey  string                `yaml:"valuesKey,omitempty"`
	TargetPath string                `yaml:"targetPath,omitempty"`
	Optional   bool                  `yaml:",omitempty"`
}

type ReleaseSpec struct {
	Interval    string              `yaml:",omitempty"`
	ChartRef    ReleaseChartRef     `yaml:"chartRef,omitempty"`
	Timeout     string              `yaml:",omitempty"`
	ReleaseName string              `yaml:"releaseName,omitempty"`
	Install     ReleaseInstall      `yaml:",omitempty"`
	Values      map[string]any      `yaml:",omitempty"`
	ValuesFrom  []ReleaseValuesFrom `yaml:"valuesFrom,omitempty"`
}

type Release struct {
	ApiVersion string          `yaml:"apiVersion" validate:"required"`
	Kind       string          `validate:"required"`
	Metadata   meta.ObjectMeta `validate:"required"`
	Spec       ReleaseSpec     `validate:"required"`
}

func NewRelease(meta meta.ObjectMeta, spec ReleaseSpec) Release {
	return Release{
		ApiVersion: "helm.toolkit.fluxcd.io/v2",
		Kind:       "HelmRelease",
		Metadata:   meta,
		Spec:       spec,
	}
}
