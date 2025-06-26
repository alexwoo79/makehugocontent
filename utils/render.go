package utils

import (
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "模板渲染失败:"+err.Error(), 500)
		return
	}
	_ = t.Execute(w, data)
}
