package jsondiff

type JSONDiff[T any] struct {
	Operation string `json:"o"`
	Value     T      `json:"v"`
}

type StringJSONDiff JSONDiff[string]
type BoolJSONDiff JSONDiff[bool]
type NumberJSONDiff JSONDiff[int]
