package core

import (
	"web/core/response"
	"fmt"
	"github.com/wolfgarnet/REST"
)

type BaseResource struct {
	Metadata REST.Metadata `json:"metadata"`
	Extensions Extensions `json:"extensions"`
}

type Extensions struct  {
	Extensions []REST.Extension
}

func newExtensions() Extensions {
	return Extensions{make([]REST.Extension, 0)}
}


func (r *BaseResource) Identifier() string {
	return r.Metadata.Id
}

func (r *BaseResource) Parent() REST.Node {
	return r.Metadata.ParentNode
}

func (r *BaseResource) GetExtensions() []REST.Extension {
	return r.Extensions.Extensions
}

/*
func (r *BaseResource) setType(t string) {
	logger.Debug("SETTING TYPE TO %v for %v", t, r)
	r.Metadata.Type = t
}

func (r *BaseResource) setId(id string) {
	r.Metadata.Id = id
}
*/

func (r *BaseResource) GetMetadata() *REST.Metadata {
	return &r.Metadata
}

func (r BaseResource) Describe() string {
	return fmt.Sprintf("Type: %v, ID: %v, Title: %v", r.Metadata.Type, r.Metadata.Id, r.Metadata.Title)
}

func (r BaseResource) GetUrlMethod(methodName, method string) REST.UrlMethod {
	logger.Debug("GET URL METHOD FOR BASE RESOURCE: %v(%v)", methodName, method)
	if methodName == "getview" && method == "get" {
		return r.doGetView
	}

	return nil
}

func (r *BaseResource) doGetView(context *REST.Context) (response.Response, error) {

	context.Request.ParseForm()
	view := context.Request.Form.Get("view")
	if len(view) == 0 {
		// return

	}

	logger.Debug("Getting view, %v, for %v", view, r.GetMetadata().Title)

	content := context.System.RenderObject(r, context)
	response := response.NewBufferedResponse()
	response.Append(content)
	return response, nil
}

/*
func (r BaseResource) String() string {
	return "RESOURCE: " + r.Title
}
*/

/*
func NewResource(id, name string) *BaseResource {
	return &BaseResource{id, name, nil}
}
*/

func (r BaseResource) UrlName() string {
	return "URL";
}

func (r *BaseResource) GetChild(name string) (REST.Node, error) {
	return nil, nil;
}

func (r BaseResource) GetDisplayName() string {
	return r.Metadata.Title
}
