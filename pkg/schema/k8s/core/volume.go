package core

type VolumeConfigMapItem struct {
	Key  string
	Path string
}

type ConfigMapVolumeSource struct {
	Name  string
	Items []VolumeConfigMapItem
}

type PVCVolumeSource struct {
	ClaimName string `yaml:"claimName,"`
}

type Volume struct {
	Name                  string                `yaml:",omitempty"`
	ConfigMap             ConfigMapVolumeSource `yaml:"configMap,omitempty"`
	PersistentVolumeClaim PVCVolumeSource       `yaml:"persistentVolumeClaim,omitempty"`
}
