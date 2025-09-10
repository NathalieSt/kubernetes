package core

type VolumeConfigMapItem struct {
	Key  string
	Path string
}

type ConfigMapVolumeSource struct {
	Name  string                `yaml:"name,"`
	Items []VolumeConfigMapItem `yaml:"items,omitempty"`
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

type Volume struct {
	Name                  string                `yaml:",omitempty"`
	ConfigMap             ConfigMapVolumeSource `yaml:"configMap,omitempty"`
	PersistentVolumeClaim PVCVolumeSource       `yaml:"persistentVolumeClaim,omitempty"`
	EmptyDir              EmptyDirVolumeSource  `yaml:"emptyDir,omitempty"`
}
