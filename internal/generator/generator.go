package generator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/nasjp/tiny-gotests/internal/model"
)

// https://github.com/golang/go/wiki/CodeReviewComments#initialisms
var initialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"CSV":   true,
	"DB":    true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"KYC":   true,
	"LHS":   true,
	"NTP":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"UUID":  true,
	"VM":    true,
	"XML":   true,
}

var tmpl = `
func Test%s%s(t *testing.T) {
  t.Pararel()

  tests := []struct{
    name string
  }

  for _, tt := range tests {
    tt := tt

    t.Run(tt.name, func(t *testing.T){

    })
  }
}
`

func Generate(fs model.Funcs) string {
	var txt string
	for _, f := range fs {
		txt += fmt.Sprintf(tmpl, capitalise(f.Receiver), capitalise(f.Name))
	}

	return txt
}

func capitalise(raw string) string {
	if raw == "" {
		return ""
	}

	if unicode.IsUpper(rune(raw[0])) {
		return raw
	}

	capitalised := []rune(raw)

	for i, r := range capitalised {
		if !unicode.IsUpper(r) {
			continue
		}

		firstSection := capitalised[0:i]

		if uppered := strings.ToUpper(string(firstSection)); initialisms[uppered] {
			return uppered + string(capitalised[i:])
		}

		return string(unicode.ToUpper(capitalised[0])) + string(capitalised[1:])
	}

	if uppered := strings.ToUpper(raw); initialisms[uppered] {
		return uppered
	}

	return string(unicode.ToUpper(capitalised[0])) + string(capitalised[1:])
}
