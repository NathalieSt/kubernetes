package meta

type MatchExpressionOperator string

const (
	In    MatchExpressionOperator = "In"
	NotIn MatchExpressionOperator = "NotIn"
)

type MatchExpression struct {
	Key      string                  `yaml:"key,omitempty"`
	Operator MatchExpressionOperator `yaml:"operator,omitempty"`
	Values   []string                `yaml:"values,omitempty"`
}

type LabelSelector struct {
	MatchLabels     map[string]string `yaml:"matchLabels,omitempty"`
	MachExpressions []MatchExpression `yaml:"matchExpressions,omitempty"`
}
