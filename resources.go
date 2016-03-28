package core

import "github.com/wolfgarnet/REST"

type resources struct {
	parent REST.Node
}

func newResources(parent REST.Node) *resources {
	return &resources{parent}
}

func (r *resources) UrlName() string {
	return "resources"
}

func (r *resources) Parent() REST.Node {
	return r.parent
}

func (r *resources) GetChild(name string) (REST.Node, error) {
	return nil, nil;
}

func (r resources) String() string {
	return "Resources"
}

func (r resources) DoGet() {

}

func (r resources) GetMetadata() *REST.Metadata {
	return nil
}

func (r resources) Identifier() string {
	return ""
}
