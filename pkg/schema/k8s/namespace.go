package k8s

/*
apiVersion: "v1" = "v1"
kind: "Namespace" = "Namespace"
metadata?: ObjectMeta
spec?: NamespaceSpec
*/
type Namespace struct {
	ApiVersion string `yaml:"apiVersion,omitempty" validate:"required"`
	Kind       string `yaml:"kind,omitempty" validate:"required"`
	Metadata   Metadata
}
