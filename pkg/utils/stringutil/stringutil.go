package stringutil

import "unicode/utf8"

func Diff(base, exclude []string) (result []string) {
	excludeMap := make(map[string]bool)
	for _, s := range exclude {
		excludeMap[s] = true
	}
	for _, s := range base {
		if !excludeMap[s] {
			result = append(result, s)
		}
	}

	return result
}

func FindString(base []string, target string) int {
	for idx, s := range base {
		if target == s {
			return idx
		}
	}

	return -1
}

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)

	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}

	return string(buf)
}
