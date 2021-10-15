package common

// Used to reserve a resource for a specific workload.
type Taint struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TolerationOperator string

const (
	// Toleration will match any value.
	TolerationOperator_Exists = "exists"
	// Toleration will only match if the value field equals the value of the taint.
	TolerationOperator_Equal = "equal"
)

type Toleration struct {
	// How the taint's value is handled.
	Operator TolerationOperator `json:"operator"`
	// The key of the taint to tolerate.
	Key string `json:"key"`
	// The value to match against the taint's value.
	// This field is ignored if 'operator' is set to 'Exists'.
	Value string `json:"value"`
	// If true, toleration must match a taint. If false, taint may be present
	// but is not required.
	Required bool `json:"required"`
}
