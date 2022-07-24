package jsondiff

import (
	"net/url"
	"strconv"
	"strings"
)

func (jsondiff *StringJSONDiff) Apply(src string) string {
	if jsondiff.Operation != OP_DMP {
		// only able to apply DMP string diffs
		return ""
	}

	diffs := strings.Split(jsondiff.Value, "\t")
	startIndex := 0
	for _, diff := range diffs {
		src, startIndex = applyDiff(diff, src, startIndex)
	}

	return src
}

func applyDiff(diff string, src string, startIndex int) (string, int) {
	operation := string(diff[0])
	newIndex := startIndex
	end := ""

	switch operation {
	case OP_INSERT:
		value, _ := url.QueryUnescape(diff[1:])
		end = src[:newIndex] + value + src[newIndex:]
		newIndex = newIndex + len(value)
	case OP_DELETE:
		value, _ := strconv.Atoi(diff[1:])
		end = src[:newIndex] + src[newIndex+value:]
	case OP_EQUAL:
		fallthrough
	default:
		value, _ := strconv.Atoi(diff[1:])
		newIndex = newIndex + value
		end = src
	}

	return end, newIndex
}
