package gotmpl

import (
	"bytes"
	"strings"
	"text/template"
)

// Render evaluates a template string with the given data.
func Render(data any, lines ...string) (out string, err error) {
	var t *template.Template

	if t, err = template.New("").Parse(strings.Join(lines, "\n")); err != nil {
		return
	}

	buf := &bytes.Buffer{}

	if err = t.Execute(buf, data); err != nil {
		return
	}

	out = buf.String()

	return
}
