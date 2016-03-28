package core

import (
	"fmt"
	"web/db"
	"github.com/wolfgarnet/templates"

	"github.com/wolfgarnet/logging"
	"web/core/registry"
	"github.com/wolfgarnet/REST"
)

var logger logging.Logger

type Configuration struct {
	TemplatePath string
	Theme string
}

type System struct {
	TemplateManager *templates.Manager
	Configuration Configuration
	db db.Collection
}

var descriptorsByUrlname map[string]Descriptor = make(map[string]Descriptor)


/*
func init() {
	logger.Info("Initializing system")
	descriptorsByUrlname = make(map[string]Descriptor)
}
*/

func (system *System) GetResource(id string) (REST.Resource, error) {
	logger.Debug("GET RESOURCE %v", id)

	object, err := system.db.Get(id)
	logger.Debug("I GOT; %v", object)
	if err != nil {
		return nil, err
	}

	resource, ok := object.(REST.Resource)
	if !ok {
		return nil, fmt.Errorf("Not a valid resource, %v", id)
	}

	/*
	br, ok := resource.(*BaseResource)
	logger.Warning("CHECK THIS OUT: %v", ok)
	if ok {
		br.outer = resource
	}
	*/

	logger.Debug("I HAVE RESOURCE: %v", resource)
	return resource, nil
}

/*
func (sys *System) AddResourceFactory(name string, factory interface{}) {
	descriptors[name] = factory
}
*/

func AddDescriptor(descriptor Descriptor) {
	logger.Info("Adding %v", descriptor)
	descriptorsByUrlname[descriptor.UrlName()] = descriptor
	//descriptorsByType.Add()
}

func AddExtension(extension interface{}) {
	registry.AddExtension(extension)
}

// Get extensions given their function name
func (sys *System) GetExtensions(fieldName string) (r []interface{}) {
	return registry.GetExtensions(fieldName)
}

func (sys *System) GetDescriptorByUrlname(name string) (Descriptor, error) {
	return descriptorsByUrlname[name], nil
}

func (sys *System) CreateResource(t string) *REST.Node {
	//val, ok := types[type]
	ok := true
	if ok {
		fmt.Println(t, "found")
	} else {
		fmt.Println(t, "NOT found")
	}
	
	return nil
}

func (sys *System) GetDescriptorNames() []Descriptor {
	l := make([]Descriptor, 0, len(descriptorsByUrlname))
	for _, v := range descriptorsByUrlname {
		l = append(l, v)
	}

	return l
}

// GetUrlPath returns the url path for a given Node
func GetUrlPath(node REST.Node) string {
	path := node.UrlName()
	for node.Parent() != nil {
		node = node.Parent()
		path = node.UrlName() + "/" + path
	}

	return path
}

func (system *System) Save(resource REST.Resource) error {
	logger.Debug("------> %v", resource)
	logger.Debug("Saving %v with id %v", resource, resource.Identifier())
	return system.db.Put(resource.Identifier(), resource)
}

func (system *System) GetId(t string) string {
	return system.db.GetId(t)
}

func (system *System) Check() string {
	return "SANITY CHECK"
}

var system *System