package httpserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// Renderer is an interface to wrap and support rendering of http responses
type Renderer interface {
	Error(w http.ResponseWriter, name string, err error)
	OK(w http.ResponseWriter, name string, data interface{})
	Status(s int, w http.ResponseWriter, name string, data interface{})
	JSON(rw http.ResponseWriter, data interface{})
	JSONStatus(status int, rw http.ResponseWriter, data interface{})
}

// TemplateProvider allows templates to be provided
type TemplateProvider interface {
	Template() string
}

type renderer struct {
	*template.Template
}

// NewRenderer provides an implementation of the Renderer interface
func NewRenderer(templateProvider []TemplateProvider) (Renderer, error) {
	funcMap := template.FuncMap{
		// The name "add" is what the function will be called in the template text.
		"add": func(a int64, b int64) int64 {
			return a + b
		},
		"sub": func(a int64, b int64) int64 {
			return a - b
		},
	}
	t := template.New("base").Funcs(funcMap)
	var err error
	for _, v := range templateProvider {
		t, err = t.Parse(v.Template())
		if err != nil {
			return nil, err
		}
	}
	return renderer{Template: t}, nil
}

func (r renderer) Error(w http.ResponseWriter, name string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	e := r.ExecuteTemplate(w, name, err)
	if e != nil {
		fmt.Println(e)
	}
}

func (r renderer) OK(w http.ResponseWriter, name string, data interface{}) {
	err := r.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (r renderer) Status(s int, w http.ResponseWriter, name string, data interface{}) {
	w.WriteHeader(s)
	err := r.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (r renderer) JSON(rw http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(js)
	if err != nil {
		fmt.Println(err)
	}
}

func (r renderer) JSONStatus(status int, rw http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, err = rw.Write(js)
	if err != nil {
		fmt.Println(err)
	}
}
