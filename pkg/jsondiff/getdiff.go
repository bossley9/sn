package jsondiff

// loosely following Google engineer Neil Fraser's implementation
// https://raw.githubusercontent.com/google/diff-match-patch/master/javascript/diff_match_patch_uncompressed.js

import (
	"strconv"
	"strings"
)

func GetDiff(s1 string, s2 string) StringJSONDiff {
	diff := StringJSONDiff{
		Operation: "d",
	}

	if s1 == s2 {
		// speedup (equality)
		diff.Value = ""
		return diff
	}
	if len(s1) == 0 {
		// speedup (complete insertion)
		diff.Value = "+" + s1
		return diff
	}
	if len(s2) == 0 {
		// speedup (complete deletion)
		diff.Value = "-" + strconv.Itoa(len(s2))
		return diff
	}

	operations := make([]string, 0)

	// speedup (common prefix/suffix)
	d1, d2, prefix, suffix := trimCommon(s1, s2)

	if len(prefix) > 0 {
		prefixOp := "=" + strconv.Itoa(len(prefix))
		operations = append(operations, prefixOp)
	}

	dmpOperations := computeDiff(d1, d2)
	for _, dmpDiff := range dmpOperations {
		operations = append(operations, dmpDiff)
	}

	if len(suffix) > 0 {
		suffixOp := "=" + strconv.Itoa(len(suffix))
		operations = append(operations, suffixOp)
	}

	diff.Value = strings.Join(operations, "\t")
	return diff
}

// given strings s1 and s2, trims common prefixes and suffixes and returns all strings
func trimCommon(s1 string, s2 string) (string, string, string, string) {
	if s1 == s2 {
		// speedup (equality)
		return s1, s2, "", ""
	}

	s1len := len(s1)
	s2len := len(s2)

	minLen := s1len
	if s2len < s1len {
		minLen = s2len
	}

	prefixIndex := 0
	for ; prefixIndex < minLen; prefixIndex++ {
		if s1[prefixIndex] != s2[prefixIndex] {
			break
		}
	}

	suffixLen := 0
	// minLen-prefixIndex prevents overlap
	for ; suffixLen < minLen-prefixIndex; suffixLen++ {
		if s1[s1len-1-suffixLen] != s2[s2len-1-suffixLen] {
			break
		}
	}

	return s1[prefixIndex : s1len-suffixLen],
		s2[prefixIndex : s2len-suffixLen],
		s1[:prefixIndex],
		s1[s1len-suffixLen:]
}

func computeDiff(s1 string, s2 string) []string {
	diffs := make([]string, 0)

	s1len := len(s1)
	s2len := len(s2)

	if s1len == 0 {
		// speedup (insertion)
		diffs = append(diffs, "+"+s2)
		return diffs
	}

	if s2len == 0 {
		// speedup (deletion)
		diffs = append(diffs, "-"+strconv.Itoa(s1len))
		return diffs
	}

	s1Index := strings.Index(s1, s2)
	if s1Index >= 0 {
		// speedup (s2 within s1)
		diffs = append(diffs, "-"+strconv.Itoa(len(s1[:s1Index])))
		diffs = append(diffs, "="+strconv.Itoa(s2len))
		diffs = append(diffs, "-"+strconv.Itoa(len(s1[s1Index+s2len:])))
		return diffs
	}

	s2Index := strings.Index(s2, s1)
	if s2Index >= 0 {
		// speedup (s1 within s2)
		diffs = append(diffs, "+"+s2[:s2Index])
		diffs = append(diffs, "="+strconv.Itoa(s1len))
		diffs = append(diffs, "+"+s2[s2Index+s1len:])
		return diffs
	}

	if s1len == 1 {
		// speedup (single character - after substring index speedup)
		diffs = append(diffs, "-1")
		diffs = append(diffs, "+"+s2)
		return diffs
	}

	if s2len == 1 {
		// speedup (single character - after substring index speedup)
		diffs = append(diffs, "-"+strconv.Itoa(s1len))
		diffs = append(diffs, "+1")
		return diffs
	}

	// future: half match: find a substring shared by both strings at least half the length of the bigger text

	// future: line mode (half match wasn't enough): compare lines of text if diff is really large (>100 lines)

	// future: bisect (adhering to Myers 1986 diff algorithm)

	// TEMP: swap and replace (inefficient - but I can't be bothered right now)
	swapDiffs := swapAndReplace(s1, s2)
	for _, diff := range swapDiffs {
		diffs = append(diffs, diff)
	}

	return diffs
}

func swapAndReplace(s1 string, s2 string) []string {
	deleteDiff := "-" + strconv.Itoa(len(s1))
	insertDiff := "+" + s2
	return []string{deleteDiff, insertDiff}
}
