package core

import (
	"web/core/response"
	"web/db"
	"regexp"
	"encoding/json"
	"errors"
	"github.com/wolfgarnet/REST"
)

type searchAction struct {
	system *System
	parent REST.Node
}

func newSearchAction(system *System, parent REST.Node) *searchAction {
	return &searchAction{system, parent}
}

func (sa *searchAction) UrlName() string {
	return "search"
}

func (sa *searchAction) Parent() REST.Node {
	return sa.parent
}

func (sa searchAction) GetUrlMethod(methodName, method string) REST.UrlMethod {
	logger.Debug("Search url method: %v", methodName)

	if methodName == "search" && method == "get" {
		return sa.doSearch
	}

	return nil
}

func (sa searchAction) GetMetadata() *REST.Metadata {
	return nil
}

func (sa searchAction) Identifier() string {
	return ""
}


func (sa searchAction) doSearch(context *REST.Context) (response.Response, error) {
	logger.Debug("RESPONE IS::::::")
	response := ProcessSearchRequestToResponse(context)

	var r map[string]interface{}
	err := json.Unmarshal(response.Json, &r)
	if err != nil {
		return nil, errors.New("Failed to read json, " + err.Error())
	}

	logger.Debug("----1111--->%v", r["resources"])

	return response, nil
}

func ProcessSearchRequest(context *REST.Context) []interface{} {

	context.Request.ParseForm()
	queryString := context.Request.FormValue("query")

	query := ParseString(queryString)

	logger.Debug("QUERY: %v", query)

	objects, err := context.System.db.Search(query)

	if err != nil {
		return nil
	}

	logger.Debug("RECORDS: %v", objects)

	renderStuff(context, objects)

	return objects
}

func ProcessSearchRequestToResponse(context *REST.Context) *response.JsonByteResponse {
	resources := ProcessSearchRequest(context)
	return response.NewJsonByteResponse(map[string]interface{}{"resources":resources})
}


func renderStuff(context *REST.Context, resources []interface{}) {
	for i, resource := range resources {
		logger.Debug("%v: %v", i, resource)

	}
	//context.Render()
}

var regex string = "(?:([\\w-]+)\\s*:\\s*(?:\"([^\"]*)\"|([\\S]+)))|(?:\"([^\"]*)\"|([\\S]+))"

func ParseString(str string) *db.Query {
	rx := regexp.MustCompile(regex)

	query := db.MakeQuery()

	matches := rx.FindAllStringSubmatch(str, -1)

	for i, m := range matches {
		logger.Debug("%v: %v", i, m)

		switch {
		// A method
		case len(m[1]) > 0:
			logger.Debug("1:")
			// With quotes
			if len(m[2]) > 0 {
				query.Clauses = append(query.Clauses, db.Where{m[1], m[2], db.Equal})

			// Without quotes
			} else if len(m[3]) > 0 {
				query.Clauses = append(query.Clauses, db.Where{m[1], m[3], db.Equal})
			}

		// Words
		case len(m[4]) > 0:

		// Just a word
		case len(m[5]) > 0:
		}
	}

	return query
}