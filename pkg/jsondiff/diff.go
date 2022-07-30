package jsondiff

// see the following concepts:
// http://simperium.github.io/jsondiff/jsondiff-js.html
// https://neil.fraser.name/writing/diff/

type JSONDiff[T any] struct {
	Operation string `json:"o"`
	Value     T      `json:"v"`
}

type StringJSONDiff JSONDiff[string]
type BoolJSONDiff JSONDiff[bool]
type Float32JSONDiff JSONDiff[float32]
type Int64JSONDiff JSONDiff[int64]
