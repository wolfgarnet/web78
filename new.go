package core

import (
	"fmt"
	"github.com/wolfgarnet/REST"
)

type newAction struct {
	system *System
}

func (r *newAction) Parent() REST.Node {
	return nil
}

func (r *newAction) UrlName() string {
	return "new"
}

func (r *newAction) GetChild(name string) (REST.Node, error) {
	d, err := r.system.GetDescriptorByUrlname(name)
	if err != nil || d == nil {
		return nil, fmt.Errorf("Unable to find descriptor, %v", err)
	}

	return nil, nil
}

func (r newAction) String() string {
	return "New action"
}
