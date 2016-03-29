package core

import (
	"github.com/wolfgarnet/templates"
	"github.com/wolfgarnet/REST"
)

// TEMPLATE runner

type Template struct {
	Renderer *templates.Renderer
}

func (runner *Template) Run(context *REST.Context) (REST.Response, error) {
	r := REST.NewBufferedResponse()
	context.Request.ParseForm()

	runner.Renderer.AddData("system", context.System)
	runner.Renderer.AddData("context", context)
	err := runner.Renderer.Render(r.Bytes)
	return r, err
}
