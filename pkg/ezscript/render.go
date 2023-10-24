package ezscript

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"setEnv": func(k, v string) string {
		return fmt.Sprintf(`export %s="%s"`, k, v)
	},
	"addEnv": func(k string, vs ...string) string {
		return fmt.Sprintf(`export %s="%s:$%s"`, k, strings.Join(vs, ":"), k)
	},
	"filepathJoin": func(vs ...string) string {
		return filepath.Join(vs...)
	},
}

// Render evaluates a template string with the given data.
func Render(data any, lines ...string) (out string, err error) {
	var t *template.Template

	if t, err = template.New("").Funcs(funcMap).Parse(strings.Join(lines, "\n")); err != nil {
		return
	}

	buf := &bytes.Buffer{}

	if err = t.Execute(buf, data); err != nil {
		return
	}

	out = buf.String()

	return
}
