package worker

import "unicode"

func removeNonPrintableChars(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.IsPrint(r) || unicode.IsSpace(r) { // Allow printable characters and space
			result = append(result, r)
		}
	}
	return string(result)
}
