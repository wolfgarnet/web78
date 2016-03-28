package core

import (
	"web/core/response"
	"reflect"
	"web/core/registry"
	"path/filepath"
	"fmt"
	"github.com/wolfgarnet/REST"
)

// The file resource
type File struct {
	BaseResource
	REST.Metadata `json:"metadata"`
	Filename string // NEEDS BSON
	Size int64 // NEEDS BSON
	Extension string `json:"extension"`
	UploadSession string // NEEDS BSON
}

type FileListener interface {
	OnFileCreated(*REST.Context, *File)

	IsApplicable(*File) bool
}

type FileTypeProvider interface {
	GetFileTypes() []string
	GetView(view string) string
}

/*
func NewFileResource(filename, session string, size int64) *File {
	file := &File{BaseResource{
		Title:filename,
	},
		Size:size,
		UploadSession: session,
	}

	return file
}
*/

// The file descriptor
type FileDescriptor struct {
	BaseDescriptor
	fileTypes map[string]FileTypeProvider
}

func (f File) Parent() REST.Node {
	return f.Metadata.ParentNode
}

func (f File) UrlName() string {
	return f.Metadata.Id
}

func (f File) Super() interface{} {
	return f.BaseResource
}

func (f *File) Save(context *REST.Context) {
	descriptor.fireOnFileSaved(f, context)
	context.System.Save(f)
}

func (f File) Type() string {
	return "file"
}

func (f File) GetUrlMethod(methodName, method string) REST.UrlMethod {
	if methodName == "get" && method == "get" {
		return f.doGetFile
	}

	return f.BaseResource.GetUrlMethod(methodName, method)
}

func (f File) doGetFile(context *REST.Context) (response.Response, error) {
	return response.NewFileResponse(f.Filename), nil
}

func (f File) FireOnFileCreated(context *REST.Context) {
	logger.Debug("Fire on file created")

	extensions := context.System.GetExtensions("OnFileCreated")
	for _, extension := range extensions {
		fl, ok := extension.(FileListener)
		if !ok {
			logger.Warning("%v is not a FileListener!", extension)
			continue
		}

		fl.OnFileCreated(context, &f)
	}
}

func (f File) GetView(view string) string {
	logger.Debug("FILE: %v, ", f)
	ftp := descriptor.GetTypeProvider(f.Extension)
	logger.Debug("FTP: %v", ftp)
	v := ftp.GetView(view)
	logger.Debug("VIEW IS %v", v)
	return v
}

func (f File) Describe() string {
	return fmt.Sprintf("Type: %v, ID: %v, Filename: %v", f.Metadata.Type, f.Identifier(), f.Filename)
}


// FILE DESCRIPTOR

var descriptor FileDescriptor
var defaultFileProvider FileTypeProvider


func (fd FileDescriptor) GetUrlMethod(methodName, method string) REST.UrlMethod {
	if methodName == "create" && method == "post" {
		return fd.doCreate
	}

	return fd.BaseDescriptor.GetUrlMethod(methodName, method)
}

/*
func (fd *FileDescriptor) Create(session *Session, parent Node, jsonData []byte) (Node, error) {
	println("data", string(jsonData))
	var file File
	if err := json.Unmarshal(jsonData, &file); err != nil {
		println("ERROR", err.Error())
	}

	fmt.Printf("File: %v\n", file)
	return &file, nil
}
*/

/*
func (fd FileDescriptor) Create(context *Context, parent Node, jsonData []byte) (Node, error) {
	var file File
	resource, err := fd.CreateBase(context, fd, &file.BaseResource)
	if err != nil {
		return nil, err
	}
	logger.Debug("======= %v, %v", resource, reflect.TypeOf(resource))

	return file, nil
}
*/

func (fd FileDescriptor) NewInstance() REST.Resource {
	return new(File)
}

func (fd FileDescriptor) doCreate(context *REST.Context) (response.Response, error) {
/*
	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		logger.Errorf("Error while reading body, %v", err)
	}
*/
	//resource, err := fd.Create(context.Session, fd, body)

	var file File
	resource, err := fd.Create(context, fd)
	if err != nil {
		return REST.MakeErrorResponse(err.Error()), nil
	}
	logger.Debug("======= %v, %v", resource, reflect.TypeOf(resource))
	//file, _ := resource.(*File)
	//id := context.System.GetId("file")
	//file.Id = id
	//file.BaseResource.Type = "file"

	/*
	// Get upload
	filename := utils.NewFilename("user-1", "myfile.jpg", "session1")
	err = utils.HandleUpload(context.Request, "file", filename)
	logger.Debug("ERROR: %v", err)
	*/

	file.Save(context)

	logger.Debug("------->%v == %v", file, file.Id)

	r := response.NewJsonResponse(map[string]interface{}{"identifier": file.Id})
	return r, nil
}

func (fd FileDescriptor) UrlName() string {
	return "files"
}

func (fd FileDescriptor) String() string {
	return "File descriptor"
}

func (fd FileDescriptor) GetType() (reflect.Type) {
	return reflect.TypeOf((*File)(nil)).Elem()
}

func (fd FileDescriptor) Save() {
	println("saving file descriptor")
}

func (fd FileDescriptor) Super() interface{} {
	return fd.BaseDescriptor
}

func (fd FileDescriptor) GetChild(name string) (REST.Node, error) {
	return nil, nil
}

func (fd FileDescriptor) AddFileTypeProvider(ftp FileTypeProvider) {
	for id, ft := range ftp.GetFileTypes() {
		logger.Debug("[%v]: %v", id, ft)
		fd.fileTypes[ft] =  ftp
	}
}

func (fd FileDescriptor) fireOnFileSaved(file *File, context *REST.Context) {
	logger.Debug("Firing on file saved for %v", file)

	// Determine file type(extension)
	//ext := filepath.Ext(file.Filename)
	//file.Type = ext
}

func (fd FileDescriptor) GetTypeProvider(t string) FileTypeProvider {
	ftp, ok := fd.fileTypes[t]
	if !ok {
		return defaultFileProvider
	}

	return ftp
}

// File listeners
// Image

type imageFileListener struct {
}

func (ifl imageFileListener) OnFileCreated(context *REST.Context, file *File) {
	ext := filepath.Ext(file.Filename)
	logger.Errorf("HEY HO! %v", ext)
	if len(ext) > 0 {
		ext = ext[1:len(ext)]
	}
	logger.Errorf("HEY HO2! %v", ext)
	file.Extension = ext
}

// File type provider

type defaultTypeProvider struct {
}

func (itp *defaultTypeProvider) GetFileTypes() []string {
	return []string{}
}

func (ift *defaultTypeProvider) GetView(view string) string {
	return "default." + view
}

type imageTypesProvider struct {
}

func (itp *imageTypesProvider) GetFileTypes() []string {
	return []string{"jpg", "jpeg", "gif", "png", "bmp"}
}

func (ift *imageTypesProvider) GetView(view string) string {
	return "image." + view
}

// Init


func init() {
	logger.Info("Initializing file")
	descriptor = FileDescriptor{
		fileTypes:make(map[string]FileTypeProvider),
	}
	AddDescriptor(&descriptor);

	// Make file types
	descriptor.AddFileTypeProvider(new(imageTypesProvider))
	defaultFileProvider = new(defaultTypeProvider)


	ifl := new(imageFileListener)
	registry.AddExtension(ifl)
}