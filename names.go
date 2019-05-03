package json2go

import (
	"bytes"
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

	result := []rune(names[0])
	for _, s := range names {
		for i, char := range []rune(s) {
			if i >= len(result) {
				break
			}

			if result[i] != char {
				result = result[:i]
				break
			}
		}
	}

	return string(result)
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
