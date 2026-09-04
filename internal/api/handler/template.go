package handler

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
)

type RenderContext struct {
	Flashes []FlashMessage
}

type FlashSeverity uint32
const (
	FlashInfo FlashSeverity = iota
	FlashWarning
	FlashError
)

type FlashMessage struct {
	Severity FlashSeverity
	Msg string
}

func NewFlash(severity FlashSeverity, msg string) FlashMessage {
	return FlashMessage{
		Severity: severity,
		Msg: msg,
	}
}

var baseTemplate = template.New("drivers")

//go:embed templates
var templateFS embed.FS

func init() {
	baseTemplate.Funcs(TemplateHelpers)

	template.Must(baseTemplate.ParseFS(templateFS, "templates/layout/*.html"))
	template.Must(baseTemplate.ParseFS(templateFS, "templates/errors/*.html"))
	template.Must(baseTemplate.ParseFS(templateFS, "templates/components/*.html"))
}

func GlobalAddFuncs(fms ...template.FuncMap) {
	for _, fm := range fms {
		baseTemplate.Funcs(fm)
	}
}

func ExtendedTemplate(fs fs.FS, templates ...string) *template.Template {
	t := template.Must(baseTemplate.Clone())
	_ = template.Must(t.ParseFS(fs, templates...))

	return t
}

var TemplateHelpers = template.FuncMap{
	"toJSON":    toJSON,
	"inputData": inputData,
}

func toJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "	")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func inputData(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("input data must have an even number of arguments")
	}

	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key := fmt.Sprint(values[i])
		m[key] = values[i+1]
	}

	return m, nil
}
