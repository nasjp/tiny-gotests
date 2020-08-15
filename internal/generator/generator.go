package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/nasjp/tiny-gotests/internal/model"
	"golang.org/x/tools/imports"
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

var pkgTempl = `
package %s_test
`

var goFile = regexp.MustCompile(`.go`)

func Generate(srcFilePath string, fn *model.Func) ([]byte, error) {
	testFilePath, err := convertTestPath(srcFilePath)
	if err != nil {
		return nil, err
	}

	f, err := parseTestFile(testFilePath, fn.Package)
	if err != nil {
		return nil, err
	}

	for _, d := range f.Decls {
		fDecl, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if testFnName := getTestName(fn); fDecl.Name.String() == testFnName {
			return nil, fmt.Errorf("func %s is aleady exist", fn.Name)
		}
	}

	decl, err := generateFnAst(fn)
	if err != nil {
		return nil, err
	}

	f.Decls = append(f.Decls, decl)
	buf := bytes.NewBuffer([]byte{})

	if err := format.Node(buf, token.NewFileSet(), f); err != nil {
		return nil, nil
	}

	fmted, err := imports.Process(testFilePath, buf.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return fmted, nil
}

func parseTestFile(path string, pkgName string) (*ast.File, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, parser.Mode(0))
	if err != nil {
		var pathError *os.PathError
		if !errors.As(err, &pathError) {
			return nil, err
		}

		f, err = parser.ParseFile(fset, path, "package "+pkgName+"_test", parser.Mode(0))
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func getTestName(fn *model.Func) string {
	return "Test" + capitalise(fn.Receiver.Type()) + capitalise(fn.Name)
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

func generateFnAst(fn *model.Func) (ast.Decl, error) {
	buf := bytes.NewBuffer([]byte{})
	args, useArgs := getArgs(fn)
	m := map[string]string{
		"name":     getTestName(fn),
		"args":     args,
		"receiver": getReceiver(fn),
		"useArgs":  useArgs,
		"want":     getWant(fn),
		"exec":     getExec(fn),
	}

	if err := tmpl.Execute(buf, m); err != nil {
		return nil, err
	}

	f, err := parser.ParseFile(token.NewFileSet(), "", buf.String(), parser.Mode(0))
	if err != nil {
		return nil, err
	}

	return f.Decls[0], nil
}

func convertTestPath(raw string) (string, error) {
	if !goFile.MatchString(raw) {
		return "", fmt.Errorf("%s is not go file", raw)
	}

	return goFile.ReplaceAllString(raw, "_test.go"), nil
}

func getArgs(fn *model.Func) (string, string) {
	var (
		args    string
		useArgs string
	)

	for _, p := range fn.Parameters {
		if args == "" {
			args = "type args struct {"
		}

		if args != "" {
			args += "\n"
		}

		args += p.Name + " "
		if p.IsPointer {
			args += "*"
		}

		if p.Package != "" {
			args += p.Package + "."
		}

		args += p.Typ
	}

	if args != "" {
		args += "\n}"
		useArgs = "args args"
	}

	return args, useArgs
}

func getPassArgs(fn *model.Func) string {
	var passArgs string

	for _, p := range fn.Parameters {
		if passArgs != "" {
			passArgs += ", "
		}

		passArgs += "tt.args." + p.Name
	}

	return passArgs
}

func getReceiver(fn *model.Func) string {
	if fn.Receiver == nil {
		return ""
	}

	receiver := "receiver "

	if fn.Receiver.IsPointer {
		receiver += "*"
	}

	if fn.Receiver.Package != "" {
		receiver += fn.Receiver.Package + "."
	}

	return receiver + fn.Receiver.Typ
}

func getWant(fn *model.Func) string {
	var want string

	gots := getGots(fn)
	for i, g := range gots {
		if want != "" {
			want += "\n"
		}

		if g == "got" {
			want += `want `

			if fn.Results[i].IsPointer {
				want += "*"
			}

			if fn.Results[i].Package != "" {
				want += fn.Results[i].Package + "."
			}

			want += fn.Results[i].Typ

			continue
		}

		if len(g) <= 3 {
			continue
		}

		cnt, err := strconv.Atoi(g[3:])
		if err != nil {
			continue
		}

		want += fmt.Sprintf("want%d ", cnt)

		if fn.Results[i].IsPointer {
			want += "*"
		}

		if fn.Results[i].Package != "" {
			want += fn.Results[i].Package + "."
		}

		want += fn.Results[i].Typ
	}

	return want
}

func getExec(fn *model.Func) string {
	gots := getGots(fn)
	run := getRun(fn)
	passArgs := getPassArgs(fn)
	errCheck := getErrCheck(fn)
	assertions := getAssertions(fn)

	var got string
	if len(gots) != 0 {
		got = strings.Join(gots, ", ") + " :="
	}

	return got + run + "(" + passArgs + ")" + "\n" + errCheck + "\n" + assertions
}

func getErrCheck(fn *model.Func) string {
	for _, r := range fn.Results {
		if r.Typ == model.ErrStr {
			return `if err != nil {
	t.Errorf("Unexpected Error: %s", err)
	return
}`
		}
	}

	return ""
}

func getGots(fn *model.Func) []string {
	got := make([]string, 0, len(fn.Results))

	var (
		cnt   int
		first int
	)

	for _, r := range fn.Results {
		if r.Typ == model.ErrStr {
			got = append(got, "err")
			continue
		}
		cnt++

		if cnt == 1 {
			got = append(got, "got")
			first = len(got) - 1

			continue
		}

		if cnt == 2 {
			got[first] = "got1"
		}

		got = append(got, fmt.Sprintf("got%d", cnt))
	}

	return got
}

func getAssertions(fn *model.Func) string {
	var assertions string

	var fStr string

	if fn.Package != "" {
		fStr = fn.Package + "."
	}

	if fn.Receiver != nil {
		fStr += fn.Receiver.Typ + "."
	}

	fStr += fn.Name

	gots := getGots(fn)
	for _, g := range gots {
		if assertions != "" {
			assertions += "\n"
		}

		if g == "got" {
			assertions += `if diff := cmp.Diff(tt.want, got); diff != "" {
	           	t.Errorf("` + fStr + `() mismatch (-want +got):\n%s", diff)
              }`

			continue
		}

		if len(g) <= 3 {
			continue
		}

		cnt, err := strconv.Atoi(g[3:])
		if err != nil {
			continue
		}

		want, got := fmt.Sprintf("want%d", cnt), fmt.Sprintf("got%d", cnt)

		assertions += `if diff := cmp.Diff(tt.` + want + `, ` + got + `); diff != "" {
t.Errorf("` + fStr + `() mismatch (-` + want + ` +` + got + `):\n%s", diff)
}`
	}

	return assertions
}

func getRun(fn *model.Func) string {
	if fn.Receiver != nil {
		return "tt.receiver." + fn.Name
	}

	return fn.Package + "." + fn.Name
}
