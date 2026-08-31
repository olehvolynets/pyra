package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
)

var baseTemplate = template.New("drivers")

func init() {
	baseTemplate.Funcs(TemplateHelpers)

	template.Must(baseTemplate.ParseGlob("view/layout/*.html"))
	template.Must(baseTemplate.ParseGlob("view/errors/*.html"))
	template.Must(baseTemplate.ParseGlob("view/components/*.html"))
}

func GlobalAddFuncs(fms ...template.FuncMap) {
	for _, fm := range fms {
		baseTemplate.Funcs(fm)
	}
}

func ExtendedTemplate(templates ...string) *template.Template {
	t := template.Must(baseTemplate.Clone())
	_ = template.Must(t.ParseFiles(templates...))

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
