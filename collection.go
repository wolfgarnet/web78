package core

import (
	"web/core/response"
	"fmt"
	"reflect"
	"github.com/wolfgarnet/REST"
)

type Collection struct {
	BaseResource
	elements []Element
}

type Element struct {
	id string `json:"id"`
}

func (c Collection) Parent() REST.Node {
	return c.Metadata.ParentNode
}

func (c Collection) UrlName() string {
	return c.Metadata.Id
}

func (c Collection) Super() interface{} {
	return c.BaseResource
}

func (c *Collection) Save(context *REST.Context) {
	context.System.Save(c)
}

func (c Collection) Type() string {
	return "collection"
}

func (c Collection) GetUrlMethod(methodName, method string) REST.UrlMethod {
	switch {
	case methodName == "add" && method == "post":
		return c.doAdd

	case methodName == "search" && method == "get":
		return c.doSearch

	case methodName == "add" && method == "get":
		return c.doAdd

	case methodName == "fetch" && method == "get":
		return c.doFetch

	case methodName == "contains" && method == "get":

	}

	return nil
}

func (c *Collection) doContains(context *REST.Context) (response.Response, error) {
	context.Request.ParseForm()
	logger.Debug("FORM IS : %v", context.Request.Form)
	id := context.Request.Form.Get("id")

	result := map[string]interface{}{"contains": c.contains(id), "id": id}
	return response.NewJsonResponse(result), nil
}

func (c *Collection) doSearch(context *REST.Context) (response.Response, error) {
	resources := ProcessSearchRequest(context)

	for _, r := range resources {
		switch t := r.(type) {
			case REST.Node:
			logger.Debug("COLLECTION NODE: %+v", t.GetMetadata().Info)
			logger.Debug("COLLECTION NODE: %+v", t.GetMetadata())
			if c.contains(t.Identifier()) {
				t.GetMetadata().Info["incollection"] = true
			} else {
				t.GetMetadata().Info["incollection"] = false
			}
		}
	}

	return response.NewJsonByteResponse(map[string]interface{}{"resources":resources}), nil
}

func (c *Collection) doAdd(context *REST.Context) (response.Response, error) {
	context.Request.ParseForm()
	logger.Debug("FORM IS : %v", context.Request.Form)
	id := context.Request.Form.Get("id")
	if len(id) == 0 {
		return nil, fmt.Errorf("Invalid request")
	}

	_, err := context.System.GetResource(id)
	if err != nil {
		return nil, fmt.Errorf("%v does not exist, %v", id, err)
	}

	// Toggle inclusion
	if c.contains(id) {
		c.remove(id)
		c.Save(context)
		return response.NewJsonResponse(map[string]interface{}{"removed": id}), nil
	} else {
		c.add(id)
		c.Save(context)
		return response.NewJsonResponse(map[string]interface{}{"added": id}), nil
	}
}

func (c *Collection) doFetch(context *REST.Context) (response.Response, error) {
	logger.Debug("Fetching {}", c)

	offset := context.GetInt("offset", 0)

	number := context.GetInt("number", 10)

	resources := c.fetch(offset, number)

	logger.Debug("-------->{}", resources)
	res := c.getFetchResponse(context, resources)
	logger.Debug("-------->{}", res)

	return response.NewJsonByteResponse(map[string]interface{}{"results":res}), nil
}

func (c *Collection) getFetchResponse(context *REST.Context, elements []Element) (r []map[string]interface{}) {
	for i, element := range elements {
		resource, err := context.System.GetResource(element.id)
		if err != nil {
			logger.Errorf("FAILED: {}", err)
			continue
		}

		avatar := context.System.RenderObject(resource, context)

		r = append(r, map[string]interface{}{"avatar" : avatar, "counter": i, "title": resource.GetDisplayName()})
	}

	return
}

func (c *Collection) fetch(offset, number int) []Element {
	diff := len(c.elements) - offset

	if diff < 0 {
		return nil
	}

	if diff - number < 0 {
		number = diff
	}

	return c.elements[offset:offset+number]
}

func (c *Collection) FetchResources(offset, number int) (resources []REST.Resource) {
	//elements := c.elements[offset:offset+number]
	/*
	for _, _ := range elements {

	}
	*/

	return
}

func (c *Collection) contains(id string) bool {
	found := false
	c.find(id, func(idx int, element *Element){
		found = true
	})

	return found
}

func (c *Collection) add(id string) bool {
	logger.Debug("Adding %v", id)

	found := false
	c.find(id, func(idx int, element *Element){
		found = true
	})

	if !found {
		c.elements = append(c.elements, Element{id})
		return true
	}

	return false
}

func (c *Collection) remove(id string) bool {
	logger.Debug("Removing %v", id)

	removed := false
	c.find(id, func(idx int, element *Element){
		c.elements = append(c.elements[:idx], c.elements[idx+1:]...)
		removed = true
	})

	return removed
}

func (c *Collection) find(id string, mf func(int, *Element)) {
	for i, e := range c.elements {
		logger.Debug("%v: %v", i, e)
		if e.id == id {
			mf(i, &e)
		}
	}
}


// Descriptor

type CollectionDescriptor struct {
	BaseDescriptor
}

func (cd CollectionDescriptor) UrlName() string {
	return "collections"
}

func (cd CollectionDescriptor) String() string {
	return "Collection descriptor"
}

func (cd CollectionDescriptor) NewInstance() REST.Resource {
	return new(Collection)
}

func (cd CollectionDescriptor) GetUrlMethod(methodName, method string) REST.UrlMethod {
	if methodName == "create" && method == "post" {
		return cd.doCreate
	}

	return cd.BaseDescriptor.GetUrlMethod(methodName, method)
}

func (cd CollectionDescriptor) doCreate(context *REST.Context) (response.Response, error) {
	node, err := cd.Create(context, cd)
	if err != nil {
		return REST.MakeErrorResponse(err.Error()), nil
	}

	resource, ok := node.(REST.Resource)
	if !ok {
		return REST.MakeErrorResponse(err.Error()), nil
	}

	context.System.Save(resource)
	logger.Debug("------->%v == %v", resource, resource.GetMetadata().Id)

	r := response.NewJsonResponse(map[string]interface{}{"identifier": resource.GetMetadata().Id})
	return r, nil
}

func (cd CollectionDescriptor) GetChild(name string) (REST.Node, error) {
	return nil, nil
}

func (cd CollectionDescriptor) GetType() (reflect.Type) {
	return reflect.TypeOf((*Collection)(nil)).Elem()
}

func (cd CollectionDescriptor) Save() {
	println("saving collection descriptor")
}

func (cd CollectionDescriptor) Super() interface{} {
	return cd.BaseDescriptor
}

/*
func (cd CollectionDescriptor) Create(session *Session, parent Node, jsonData []byte) (Node, error) {
	println("data", string(jsonData))
	var collection Collection
	if err := json.Unmarshal(jsonData, &collection); err != nil {
		logger.Errorf("Error while unmarshalling, %v", err)
	}

	logger.Debug("Collection: %v\n", collection)
	return &collection, nil
}
*/


func init() {
	fmt.Println("HEJ")
	AddDescriptor(new(CollectionDescriptor));
}
