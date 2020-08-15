package generator

import "text/template"

var tmpl = template.Must(template.New("tmpl").Parse(`package tmp_test
func {{.name}}(t *testing.T) {
	t.Parallel()
	{{.args}}
	tests := []struct {
		name string
		{{.receiver}}
		{{.useArgs}}
		{{.want}}
	}{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			{{.exec}}
		})
	}
}`))
