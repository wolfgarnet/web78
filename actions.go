package core

import (
	"web/core/response"
	"web/core/utils"
	"fmt"
	"path/filepath"
	"reflect"
	"github.com/wolfgarnet/REST"
)


// RESOURCE ACTION

type resourceAction struct {
	system *System
	parent REST.Node
}

func newResourceAction(system *System, parent REST.Node) *resourceAction {
	return &resourceAction{system, parent}
}

func (r *resourceAction) UrlName() string {
	return "resource"
}

func (r *resourceAction) Parent() REST.Node {
	return r.parent
}

func (r *resourceAction) GetChild(name string) (REST.Node, error) {
	logger.Debug("GETTING RESOURCE %v", name)

	resource, err := r.system.GetResource(name)
	if err != nil {
		logger.Errorf("FAILED: %v", err)
		return nil, err
	}

	return resource, nil

	/*
	node, ok := resource.(Node)
	if !ok {
		return nil, fmt.Errorf("Not a valid node, %v", name)
	}

	return node
	*/
}

func (r resourceAction) String() string {
	return "Resource action"
}

func (r resourceAction) GetMetadata() *REST.Metadata {
	return nil
}

func (r resourceAction) Identifier() string {
	return ""
}


// UPLOAD ACTION

type uploadAction struct {
	parent REST.Node
}

func newUploadAction(parent REST.Node) *uploadAction {
	return &uploadAction{parent}
}

func (ua *uploadAction) UrlName() string {
	return "upload"
}

func (ua *uploadAction) Parent() REST.Node {
	return ua.parent
}

func (ua uploadAction) GetUrlMethod(methodName, method string) REST.UrlMethod {
	logger.Debug("I AM HERE")
	if methodName == "upload" && method == "post" {
		return ua.upload
	}

	return nil
}

func (ua *uploadAction) upload(context *REST.Context) (response.Response, error) {
	fileUploads, err := utils.HandleUpload(context.Request, "files", utils.NewFilenamer("user-1"))
	if err != nil {
		return nil, err
	}

	logger.Debug("RESPONSE UPLOAD: %v", fileUploads[0])
	system, ok := context.System.(*System)
	if !ok {
		panic("The system is not correct!")
	}
	descriptor, err := system.GetDescriptorByUrlname("files")
	if err != nil {
		// TODO delete file
		return nil, err
	}
	if descriptor == nil {
		return nil, fmt.Errorf("Descriptor not found")
	}

	logger.Debug("DESCRIPTOR: %v", descriptor)

	// Create file resources
	for _, fu := range fileUploads {

		//node, nil := descriptor.Create(context, descriptor)
		resource := descriptor.NewInstance()
		logger.Debug("NODE: %v -- %v", resource, reflect.TypeOf(resource))
		resource, err = FillBase(resource, system.GetId(resource.Type()))
		if err != nil {
			logger.Errorf("FAILED: %v", err)
			return REST.MakeErrorResponse("WAAAG, err: " + err.Error()), nil
		}
		logger.Debug("STRUCT: %+v", resource)

		//node, err := descriptor.Create(context.Session, descriptor, []byte("{}"))
		if err != nil {
			logger.Warning("Unable to create file, %v", err)
		}

		file, ok := resource.(*File)
		if !ok {
			logger.Warning("Not a file resouce")
			continue
		}

		file.Size = fu.Size
		file.Title = fu.Name
		file.Filename = fu.Filename
		file.UploadSession = fu.UploadSession
		file.Id = system.GetId("file")

		fu.Rid = file.Id
		fu.Url = GetUrlPath(file)

		ext := filepath.Ext(file.Filename)
		if len(ext) > 0 {
			ext = ext[1:len(ext)]
		}
		file.Extension = ext

		// Fire the listeners just before the file is saved
		file.FireOnFileCreated(context)

		file.Save(context)

		logger.Debug("RESOURCE: %v", file.Describe())
		logger.Debug("STRUCT: %+v", file)
	}

	return response.NewJsonResponse(map[string]interface{}{"files":fileUploads}), nil
}

func (ua uploadAction) String() string {
	return "Upload action"
}


func (ua uploadAction) GetMetadata() *REST.Metadata {
	return nil
}

func (ua uploadAction) Identifier() string {
	return ""
}
