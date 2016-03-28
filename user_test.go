package core

import (
	"testing"
	"fmt"
	"encoding/json"
)

func TestCreateUser(t *testing.T) {

}

func TestTypeImplements(t *testing.T) {
	ud := UserDescriptor{}

	/*
	data := map[string]interface{}{
		"id": "user-1",
		"name": "name",
	}
	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))
	fmt.Println(data)
	*/
	//jsonData := []byte(`{"resource": {"id": "user-1", "name":"wolle"}}`)

	r := BaseResource{"user-1", "wolle", nil}
	user := User{r, "knolle"}
	json := jsonDump("USER", user)

	session := &Session{}

	//node, _ := ud.Create(session, nil, json)
	//node := ud.NewInstance()
	//ud.CreateFromJson(node, nil, )
	newuser := node.(*User)
	println("HEJJEJEJEJE", newuser)
	println("NAME;", newuser.Username(), newuser.Password)
}

func jsonDump(label string, thingy interface{}) []byte {
	serialized, _ := json.Marshal(thingy)
	fmt.Printf(label+" as JSON:\n  %s\n", serialized)
	return serialized
}