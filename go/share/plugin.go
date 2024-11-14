package share

import (
	"strings"
)

// EscapeDomain return an escaped string,
// suitable for use as an identifier with magefile scripts etc
func EscapeDomain(domain string) string {
	return strings.ReplaceAll(domain, ".", "-")
}
