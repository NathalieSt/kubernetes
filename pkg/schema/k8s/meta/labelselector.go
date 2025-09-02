package meta

type LabelSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels,omitempty"`
}
