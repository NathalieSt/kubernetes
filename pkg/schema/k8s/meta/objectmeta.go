package meta

type ObjectMeta struct {
	Annotations map[string]string `yaml:",omitempty"`
	Labels      map[string]string `yaml:",omitempty"`
	Name        string            `yaml:",omitempty"`
	Namespace   string            `yaml:",omitempty"`
}
