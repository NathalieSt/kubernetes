package vaultsecretsoperator

type Kubernetes struct {
	Audiences              []string `yaml:"audiences,omitempty"`
	Mount                  string   `yaml:"mount,omitempty"`
	Role                   string   `yaml:"role,omitempty"`
	ServiceAccount         string   `yaml:"serviceAccount,omitempty"`
	TokenExpirationSeconds int      `yaml:"tokenExpirationSeconds,omitempty"`
}
