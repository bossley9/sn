package jsondiff

import (
	"net/url"
	"strconv"
	"strings"

	l "github.com/bossley9/sn/pkg/logger"
)

func (jsondiff *StringJSONDiff) PrettyPrint(source string, content string) {
	l.PrintPlain(source + "\n")
	contentIndex := 0

	for _, diff := range strings.Split(jsondiff.Value, "\t") {
		operation := string(diff[0])
		value, _ := url.QueryUnescape(diff[1:])
		// equal diff buffer for context
		const trailingCharCount = 50

		switch operation {
		case OP_INSERT:
			l.PrintDiffInsert(value)
		case OP_DELETE:
			charCount, _ := strconv.Atoi(value)
			l.PrintDiffDelete(content[contentIndex : contentIndex+charCount])
			contentIndex = contentIndex + charCount
		case OP_EQUAL:
			fallthrough
		default:
			charCount, _ := strconv.Atoi(value)
			if charCount > trailingCharCount*2 {
				trailingStart := content[contentIndex : contentIndex+trailingCharCount]
				trailingEnd := content[contentIndex+charCount-trailingCharCount : contentIndex+charCount]
				l.PrintPlain(trailingStart + "..." + trailingEnd)
			} else {
				l.PrintPlain(content[contentIndex : contentIndex+charCount])
			}
			contentIndex = contentIndex + charCount
		}
	}

	l.PrintPlain("\n")
}
