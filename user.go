package core

import (
	"fmt"
	"encoding/json"
	"errors"
	"reflect"
	"web/core/response"
	"github.com/wolfgarnet/REST"
)

// The user resource
type User struct {
	BaseResource
	Password string `json:"password"`
}

// The user descriptor
type UserDescriptor struct {
	BaseDescriptor
}

func (u User) Username() string {
	return u.Metadata.Title;
}

type user User

func (u *User) UnmarshalJSON(b []byte) (err error) {
	fmt.Println("RES: ", string(b))
	r2 := user{}
	if err = json.Unmarshal(b, &r2); err == nil {
		*u = User(r2)
		fmt.Println("RES2: ", u.Metadata.Title)
		return nil
	}

	return errors.New("Not a valid Resource json")
}

func (u *User) GetUrlMethod(methodName, method string) REST.UrlMethod {
	logger.Debug("GET URL METHID FOR USER: %v", methodName)
	m := u.BaseResource.GetUrlMethod(methodName, method)
	logger.Debug("----->%v", m)
	return m
}

/*
func NewUser(name string) *User {
	u := &User{}
	r := NewResource("user-1", name)
	u.BaseResource = *r

	return u
}
*/

/*
func (u User) String() string {
	return "HEJ " + u.Title;
}
*/

func (u User) GetChild(name string) (REST.Node, error) {
	return nil, nil
}

func (u User) Parent() REST.Node {
	return nil
}

func (u *User) Save(context *REST.Context) {
	context.System.Save(u)
}

/*
func (u User) GetId() string {
	return u.Id
}
*/

func (u User) UrlName() string {
	return u.Metadata.Id
}

func (u User) Super() interface{} {
	return u.BaseResource
}

func (u User) Type() string {
	return "user"
}


/// USER DESCRIPTOR

/*
func (ud *UserDescriptor) Create(session *Session, parent Node, jsonData []byte) (Node, error) {
	println("data", string(jsonData))
	var user User
	if err := json.Unmarshal(jsonData, &user); err != nil {
		println("ERROR", err.Error())
	}

	fmt.Printf("USER: %v\n", user)
	return &user, nil
}
*/

func (ud UserDescriptor) UrlName() string {
	return "users"
}

func (ud UserDescriptor) String() string {
	return "User descriptor"
}

func (ud UserDescriptor) GetType() (reflect.Type) {
	return reflect.TypeOf((*User)(nil)).Elem()
}

func (ud UserDescriptor) Save() {
	println("saving user descriptor")
}

func (ud UserDescriptor) GetChild(name string) (REST.Node, error) {
	return nil, nil
}

func (ud UserDescriptor) GetUrlMethod(methodName, method string) REST.UrlMethod {
	if methodName == "create" && method == "post" {
		return ud.doCreate
	}

	return ud.BaseDescriptor.GetUrlMethod(methodName, method)
}

func (ud UserDescriptor) NewInstance() REST.Resource {
	logger.Debug("RETURNING NEW USER")
	return new(User)
}

func (ud UserDescriptor) doCreate(context *REST.Context) (response.Response, error) {
	/*
	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		logger.Errorf("Error while reading body, %v", err)
	}
	*/

	//resource, err := ud.Create(context.Session, ud, body)
	var user *User
	resource, err := ud.Create(context, ud)
	if err != nil {
		return REST.MakeErrorResponse(err.Error()), nil
	}
	//FillBase(context, ud, user)
	logger.Debug("======= %v, %v", resource, reflect.TypeOf(resource))
	//user, _ := resource.(*User)
	//id := context.System.GetId("user")
	//user.Id = id
	//user.BaseResource.Type = "user"
	context.System.Save(user)
	logger.Debug("------->%v == %v", user, user.Metadata.Id)

	r := response.NewJsonResponse(map[string]interface{}{"identifier": user.Metadata.Id})
	return r, nil
}

func (ud UserDescriptor) Super() interface{} {
	return ud.BaseDescriptor
}


// Init


func init() {
	fmt.Println("HEJ")
	AddDescriptor(new(UserDescriptor));
}