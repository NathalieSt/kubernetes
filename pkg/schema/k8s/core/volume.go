package core

type VolumeConfigMapItem struct {
	Key  string
	Path string
}

type ConfigMapVolumeSource struct {
	Name        string                `yaml:"name,"`
	Items       []VolumeConfigMapItem `yaml:"items,omitempty"`
	DefaultMode int32                 `yaml:"defaultMode,omitempty"`
}

type PVCVolumeSource struct {
	ClaimName string `yaml:"claimName,"`
}

type EmptyDirVolumeSourceMedium string

const (
	Memory EmptyDirVolumeSourceMedium = "Memory"
)

type EmptyDirVolumeSource struct {
	Medium EmptyDirVolumeSourceMedium `yaml:"medium,"`
}

type SecretVolumeItem struct {
	Key  string `yaml:"key,omitempty"`
	Path string `yaml:"path,omitempty"`
}

type SecretVolumeSource struct {
	SecretName string             `yaml:"secretName,omitempty"`
	Items      []SecretVolumeItem `yaml:"items,omitempty"`
}

type HostPathType = string

const (
	Directory HostPathType = "Directory"
)

type HostPath struct {
	Path string       `yaml:"path,omitempty"`
	Type HostPathType `yaml:"type,omitempty"`
}

type Volume struct {
	Name                  string                `yaml:",omitempty"`
	ConfigMap             ConfigMapVolumeSource `yaml:"configMap,omitempty"`
	PersistentVolumeClaim PVCVolumeSource       `yaml:"persistentVolumeClaim,omitempty"`
	EmptyDir              EmptyDirVolumeSource  `yaml:"emptyDir,omitempty"`
	Secret                SecretVolumeSource    `yaml:"secret,omitempty"`
	HostPath              HostPath              `yaml:"hostPath,omitempty"`
}
