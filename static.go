package core

import (
	"web/core/response"
	"github.com/wolfgarnet/REST"
)

type staticAction struct  {
	parent REST.Node
}

func newStaticAction(parent REST.Node) *staticAction {
	return &staticAction{parent}
}

func (r *staticAction) Parent() REST.Node {
	return r.parent
}

func (r *staticAction) UrlName() string {
	return "static"
}

func (r *staticAction) Autonomize(context *REST.Context) (response.Response, error) {
	p := context.Request.URL.String()[7:len(context.Request.URL.String())]
	path := "static" + p
	logger.Debug("P: %v, %v", p, path)

	return response.NewFileResponse(path), nil
}

func (r staticAction) String() string {
	return "Static action"
}

func (r staticAction) GetMetadata() *REST.Metadata {
	return nil
}

func (r staticAction) Identifier() string {
	return ""
}
