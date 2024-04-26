package views

import (
	"embed"
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed *html
var tmplate embed.FS

func TemplateExecuted(w http.ResponseWriter, v any, filename ...string) error {
	t, err := template.ParseFS(tmplate, filename...)
	if err != nil {
		return err
	}
	return t.Execute(w, v)
}

func TemplateMap(w http.ResponseWriter, v any, Fmap template.FuncMap, filename ...string) error {
	t:=template.New("bro").Funcs(Fmap)
	_,err:=t.ParseFS(tmplate,filename...)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w,t.Name(),v)
}
