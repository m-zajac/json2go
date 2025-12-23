package json2go

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// attrName converts json field name to pretty struct attribute name
func attrName(fieldName string) string {
	var b bytes.Buffer

	var words []string

	var i int
	for s := fieldName; s != ""; s = s[i:] { // split on upper letter or _
		i = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if i <= 0 {
			i = len(s)
		}
		word := s[:i]
		words = append(words, strings.Split(word, "_")...)
	}

	// words := strings.Split(fieldName, "_")
	for i, word := range words {
		if u := strings.ToUpper(word); commonInitialisms[u] {
			b.WriteString(u)
			continue
		}

		word = removeInvalidChars(word, i == 0) // on 0 remove first digits
		if len(word) == 0 {
			continue
		}

		out := strings.ToUpper(string(word[0]))
		if len(word) > 1 {
			out += strings.ToLower(word[1:])
		}
		b.WriteString(out)
	}

	if b.Len() == 0 { // check if this is number
		if _, err := strconv.Atoi(fieldName); err == nil {
			b.WriteString("Key")
			b.WriteString(fieldName)
		}
	}

	return b.String()
}

func removeInvalidChars(s string, removeFirstDigit bool) string {
	var buf bytes.Buffer

	for _, b := range []byte(s) {
		if b >= 97 && b <= 122 { // a-z
			buf.WriteByte(b)
			continue
		}
		if b >= 65 && b <= 90 { // A-Z
			buf.WriteByte(b)
			continue
		}
		if b >= 48 && b <= 57 { // 0-9
			if !removeFirstDigit || buf.Len() > 0 {
				buf.WriteByte(b)
				continue
			}
		}
	}

	return buf.String()
}

func extractCommonName(names ...string) string {
	if len(names) == 0 {
		return ""
	}

	// Try word-based prefix first
	var prefixParts []string
	firstWords := splitToWords(names[0])
	for i := 1; i <= len(firstWords); i++ {
		prefix := strings.Join(firstWords[:i], "_")
		allMatch := true
		for _, s := range names[1:] {
			if !strings.HasPrefix(strings.ToLower(s), prefix) {
				allMatch = false
				break
			}
		}
		if allMatch {
			prefixParts = firstWords[:i]
		} else {
			break
		}
	}

	// Try word-based suffix
	var suffixParts []string
	lastWords := splitToWords(names[0])
	for i := 1; i <= len(lastWords); i++ {
		suffix := strings.Join(lastWords[len(lastWords)-i:], "_")
		allMatch := true
		for _, s := range names[1:] {
			if !strings.HasSuffix(strings.ToLower(s), suffix) {
				allMatch = false
				break
			}
		}
		if allMatch {
			suffixParts = lastWords[len(lastWords)-i:]
		} else {
			break
		}
	}

	prefix := strings.Join(prefixParts, "_")
	suffix := strings.Join(suffixParts, "_")

	// If word-based failed, try rune-based prefix
	if len(prefix) == 0 && len(suffix) == 0 {
		resPrefix := []rune(names[0])
		for _, s := range names {
			for i, char := range []rune(s) {
				if i >= len(resPrefix) {
					break
				}
				if resPrefix[i] != char {
					resPrefix = resPrefix[:i]
					break
				}
			}
		}
		prefix = strings.Trim(string(resPrefix), "-_")
	}

	if len(prefix) >= len(suffix) {
		return prefix
	}
	return suffix
}

func splitToWords(s string) []string {
	var words []string
	var current strings.Builder
	for _, r := range s {
		if r == '_' || r == '-' || unicode.IsUpper(r) {
			if current.Len() > 0 {
				words = append(words, strings.ToLower(current.String()))
				current.Reset()
			}
			if r != '_' && r != '-' {
				current.WriteRune(r)
			}
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		words = append(words, strings.ToLower(current.String()))
	}
	return words
}

func nameFromNames(names ...string) string {
	if len(names) == 0 {
		return ""
	}

	result := names[0]
	for i, k := range names[1:] {
		result = result + "_" + k
		if i > 0 && len(result) > 3 {
			return result
		}
	}

	return result
}

func nameFromNamesCapped(names ...string) string {
	if len(names) == 0 {
		return ""
	}

	const maxParts = 3
	result := names[0]
	for i, k := range names[1:] {
		if i >= maxParts-1 {
			result = fmt.Sprintf("%sAnd%dMore", result, len(names)-maxParts+1)
			return result
		}
		result = result + "_" + k
	}

	return result
}

func nextName(name string) string {
	if name == "" {
		return "1"
	}

	re := regexp.MustCompile(`\d+$`)
	subs := re.FindStringSubmatch(name)
	if len(subs) == 0 {
		return name + "2"
	}

	num, err := strconv.Atoi(subs[0])
	if err != nil {
		return name + "2"
	}

	return re.ReplaceAllString(name, strconv.Itoa(num+1))
}

// commonInitialisms is a set of common initialisms.
//
// source: https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}
