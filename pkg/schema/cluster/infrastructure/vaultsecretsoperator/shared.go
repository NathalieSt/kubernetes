package vaultsecretsoperator

type Kubernetes struct {
	Audiences              []string `yaml:"audiences"`
	Mount                  string   `yaml:"mount"`
	Role                   string   `yaml:"role"`
	ServiceAccount         string   `yaml:"serviceAccount"`
	TokenExpirationSeconds int      `yaml:"tokenExpirationSeconds"`
}
