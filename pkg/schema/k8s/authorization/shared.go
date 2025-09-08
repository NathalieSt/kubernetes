package authorization

type Rule struct {
	APIGroups       []string `yaml:"apiGroups,omitempty"`
	Resources       []string `yaml:"resources,omitempty"`
	Verbs           []string `yaml:"verbs,omitempty"`
	NonResourceURLs []string `yaml:"nonResourceURLs,omitempty"`
}

type RoleRef struct {
	APIGroup string `yaml:"apiGroup,omitempty"`
	Kind     string `yaml:"kind,omitempty"`
	Name     string `yaml:"name,omitempty"`
}
type Subject struct {
	Kind      string `yaml:"kind,omitempty"`
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}
