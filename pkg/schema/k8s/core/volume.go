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
	Name string
}

type Volume struct {
	Name                  string
	ConfigMap             *ConfigMapVolumeSource `yaml:",omitempty"`
	PersistentVolumeClaim *PVCVolumeSource       `yaml:",omitempty"`
}
