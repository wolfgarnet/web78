package core

import (
	"fmt"
	
	"errors"
	"github.com/wolfgarnet/REST"
)

type root struct {
	children map[string]REST.Node
	system *System
}

func (r *root) Parent() REST.Node {
	return nil
}

func (r *root) UrlName() string {
	return ""
}

func (r *root) GetChild(name string) (REST.Node, error) {
	println("Child of root", name)
	v, ok := r.children[name]
	if ok {
		return v, nil
	}

	d, err := r.system.GetDescriptorByUrlname(name)
	if err == nil && d != nil {
		return d, nil
	}

	return nil, errors.New(name + " not found")
}

func (r root) AddChild(name string, child REST.Node) {
	fmt.Println("Adding ", child)
	r.children[name] = child
}

func NewRoot(system *System) *root {
	root := root{make(map[string]REST.Node), system}
	return &root
}

func (r root) String() string {
	return "Root"
}

func (r root) GetMetadata() *REST.Metadata {
	return nil
}

func (r root) Identifier() string {
	return ""
}
