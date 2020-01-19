package template

import (
	"io"

	p "github.com/majori/goco/pkg/parser"
)

type Template struct {
	Source *string
}

func New(source *string) *Template {
	return &Template{source}
}

func (tmpl Template) PreProcess(state *p.States) {

}

func (tmpl Template) Execute(writer io.Writer, settings *p.Settings) {

}
