package extensions

import (
	"web/core"
	"web/core/registry"
	"github.com/wolfgarnet/logging"
)

var logger logging.Logger

type getable interface {
	ProcessGetable() interface{}
}

// Action,

type GetAction struct {

}

func (g GetAction) ProcessAction() {

}

func (g GetAction) GetName() string {
	return "get"
}

func (g GetAction) IsApplicable(node core.Node) bool {
	return true
}

func (g GetAction) GetChild(name string) (core.Node, error) {
	return nil, nil
}

func (g GetAction) Identifier() string {
	return "GetAction"
}

func (g GetAction) UrlName() string {
	return "get"
}

func (g GetAction) Parent() core.Node {
	return nil
}

func (g GetAction) GetMetadata() *core.Metadata {
	return nil
}


func (g GetAction) String() string {
	return "Getable"
}

func init() {
	logger.Info("Initializing getable")
	action := GetAction{}
	registry.AddExtension(action);
}