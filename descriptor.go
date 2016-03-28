package core

import (
	"web/core/response"
	"reflect"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/wolfgarnet/REST"
)

type Descriptor interface {
	REST.Node

	Create(context *REST.Context, descriptor Descriptor) (REST.Node, error)
	GetChild(name string) (REST.Node, error)
	Save()
	GetType() reflect.Type
	NewInstance() REST.Resource
	//FillBase(context *Context, resource Resource) (Resource, error)
}

type BaseDescriptor struct  {
	//outer Resource
}

/*
func (d BaseDescriptor) NewInstance(session Session, parent *Node, jsonData []byte) (Node, error) {
	instance, err := d.Create(session, parent, jsonData)
	return instance, err
}
*/

func (d BaseDescriptor) GetUrlMethod(methodName, method string) REST.UrlMethod {
	logger.Debug("I AM: %v", d)
	if methodName == "create" && method == "post" {
		logger.Debug("Returning create FAIL")
		return func(context *REST.Context) (response.Response, error) {
			return REST.MakeErrorResponse("NO CREATE METHOD CREATED!!!!!"), nil
		}
	}

	return nil
}

func (d BaseDescriptor) GetMetadata() *REST.Metadata {
	return nil
}

func (d BaseDescriptor) Identifier() string {
	return ""
}

/*
func (d BaseDescriptor) doCreate(context *Context) (response.Response, error) {
	r := response.NewJsonResponse(map[string]interface{}{"identifier":"user-1"})
	return r, nil
}
*/

/*
func (d BaseDescriptor) Create(context *Context, parent Node) (Node, error) {
	return nil, nil
}
*/

func (d BaseDescriptor) CreateFromJson(resource REST.Resource, parent REST.Node, jsonData []byte) (REST.Resource, error) {
	logger.Debug("JSON: %v", string(jsonData))
	if err := json.Unmarshal(jsonData, resource); err != nil {
		println("ERROR", err.Error())
		return nil, err
	}

	logger.Debug("RESULT: %v", resource)

	return resource, nil
}

func (d BaseDescriptor) Create(context *REST.Context, descriptor Descriptor) (REST.Node, error) {
	logger.Debug("CREATING NODE: %v", descriptor)

	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		logger.Errorf("Error while reading body, %v", err)
		return nil, err
	}

	resource := descriptor.NewInstance()
	logger.Debug("BOOOOOOOOOM!")
	resource, err = d.CreateFromJson(resource, d, body)
	if err != nil {
		return nil, fmt.Errorf("Unable to create from json, %v", err.Error())
	}

	logger.Debug("-------->%v", reflect.TypeOf(resource))
	logger.Debug("-------->%v", resource.GetMetadata())

	return FillBase(resource, context.System.GetId(resource.Type()))
}

func FillBase(resource REST.Resource, id string) (REST.Resource, error) {
	resource.GetMetadata().Type = resource.Type()
	resource.GetMetadata().Id = id

	logger.Debug("RESOURCE: %v\n", resource)

	resource.GetMetadata().Info = make(map[string]interface{})
	resource.GetMetadata().Info["test"] = "created"

	return resource, nil
}

/*
func (d BaseDescriptor) FillBase(context *Context, resource Resource) (Resource, error) {
	//resource.g
	logger.Debug("DESCRIPTOR: %v", d)
	logger.Debug(":::::: %v -> %v", d.GetType(), resource)
	t := strings.ToLower(d.GetType().Name())
	logger.Debug("THE TYPE: %v - %v", descriptor.GetType(), t)
	resource.GetMetadata().Type = t
	resource.GetMetadata().Id = context.System.GetId(t)

	logger.Debug("RESOURCE: %v\n", resource)


	return resource, nil
}
*/

/*
func (d BaseDescriptor) Create(session Session, parent *Node, jsonData []byte) (Node, error) {
	return nil, nil
}
*/

func (d BaseDescriptor) GetType() reflect.Type {
	return nil
}

func (d BaseDescriptor) UrlName() string {
	return ""
}

func (d BaseDescriptor) GetDisplayName() string {
	return ""
}

func (d BaseDescriptor) GetId() string {
	return ""
}

func (d BaseDescriptor) GetJsonId() string {
	return ""
}

func (d BaseDescriptor) Parent() REST.Node {
	return nil
}
/*
func (d BaseDescriptor) GetChild(name string) (Node, error) {
	return nil, nil
}

func (d BaseDescriptor) GetType() reflect.Type {
	return nil
}

func (d BaseDescriptor) NewInstance() Resource {
	return nil
}

func (d BaseDescriptor) Save() {

}
*/