package k8s

/*
schema ObjectMeta:
    annotations?: {str:str}
    creationTimestamp?: str
    deletionGracePeriodSeconds?: int
    deletionTimestamp?: str
    finalizers?: [str]
    generateName?: str
    generation?: int
    labels?: {str:str}
    managedFields?: [ManagedFieldsEntry]
    name?: str
    namespace?: str
    ownerReferences?: [OwnerReference]
    resourceVersion?: str
    selfLink?: str
    uid?: str

*/

type Metadata struct {
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Name        string            `yaml:"name,omitempty"`
	Namespace   string            `yaml:"namespace,omitempty"`
}
