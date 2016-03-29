package core

import (
	"web/core/response"
	"log"
	"github.com/wolfgarnet/REST"
)

type themeAction struct  {
	parent REST.Node
}

func newThemeAction(parent REST.Node) *themeAction {
	return &themeAction{parent}
}

func (r *themeAction) Parent() REST.Node {
	return r.parent
}

func (r *themeAction) UrlName() string {
	return "theme"
}

func (r *themeAction) Autonomize(context *REST.Context) (response.Response, error)  {
	//log.Printf("--->%v", context.Request.URL)
	p := context.Request.URL.String()[6:len(context.Request.URL.String())]
	theme := context.Session.Get("theme", "default")
	platform := context.Session.Get("platform", "desktop")
	path := "themes/" + theme + "/" + platform + p
	log.Printf("P: %v, %v", p, path)
	//path :=
	//return response.NewFileResponse("themes/default/desktop/style.css")
	return response.NewFileResponse(path), nil
}

func (r themeAction) String() string {
	return "Theme action"
}

func (r themeAction) GetMetadata() *REST.Metadata {
	return nil
}

func (r themeAction) Identifier() string {
	return ""
}
