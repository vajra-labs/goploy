package apidoc

import (
	"fmt"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type ResponseEntry struct {
	status   string
	response *huma.Response
}

// Responses builds a map[string]*huma.Response from entries.
// Same as TypeScript: responses: { [HttpStatus.OK]: jsonContent(...), ... }
func Responses(entries ...ResponseEntry) map[string]*huma.Response {
	m := map[string]*huma.Response{}
	for _, e := range entries {
		m[e.status] = e.response
	}
	return m
}

// JsonContent creates a JSON response entry.
func JsonContent(
	api huma.API,
	status int,
	v any,
	description string,
) ResponseEntry {
	return ResponseEntry{
		status: fmt.Sprintf("%d", status),
		response: &huma.Response{
			Description: description,
			Content: map[string]*huma.MediaType{
				"application/json": {Schema: schemaFor(api, v)},
			},
		},
	}
}

// TextContent creates a text/plain response entry.
func TextContent(status int, description string) ResponseEntry {
	return ResponseEntry{
		status: fmt.Sprintf("%d", status),
		response: &huma.Response{
			Description: description,
			Content: map[string]*huma.MediaType{
				"text/plain": {Schema: &huma.Schema{Type: "string"}},
			},
		},
	}
}

// ErrContent creates an error response entry using custom HttpError schema.
func ErrContent(status int, description string) ResponseEntry {
	return ResponseEntry{
		status: fmt.Sprintf("%d", status),
		response: &huma.Response{
			Description: description,
			Content: map[string]*huma.MediaType{
				"application/json": {
					Schema: &huma.Schema{Ref: "#/components/schemas/HttpError"},
				},
			},
		},
	}
}

// Body creates a request body schema.
func Body(
	api huma.API,
	v any,
	required bool,
	description string,
) *huma.RequestBody {
	return &huma.RequestBody{
		Required:    required,
		Description: description,
		Content: map[string]*huma.MediaType{
			"application/json": {
				Schema: schemaFor(api, v),
			},
		},
	}
}

// QueryParams generates query parameters from a struct's fields.
// Uses `json` tag as param name and `doc` tag as description.
// Uses `validate:"required"` to mark required params.
func QueryParams(v any) []*huma.Param {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var params []*huma.Param
	for i := range t.NumField() {
		field := t.Field(i)
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			continue
		}
		description := field.Tag.Get("doc")
		required := false
		if v := field.Tag.Get("validate"); v != "" {
			for _, rule := range splitComma(v) {
				if rule == "required" {
					required = true
					break
				}
			}
		}
		schema := fieldSchema(field.Type)
		params = append(params, &huma.Param{
			Name:        name,
			In:          "query",
			Description: description,
			Required:    required,
			Schema:      schema,
		})
	}
	return params
}

// PathParams generates path parameters from a struct's fields.
// Uses `json` tag as param name and `doc` tag as description.
func PathParams(v any) []*huma.Param {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var params []*huma.Param
	for i := range t.NumField() {
		field := t.Field(i)
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			continue
		}
		description := field.Tag.Get("doc")
		params = append(params, &huma.Param{
			Name:        name,
			In:          "path",
			Description: description,
			Required:    true,
			Schema:      fieldSchema(field.Type),
		})
	}
	return params
}

// schemaFor generates a huma.Schema from any Go struct using the API registry.
func schemaFor(api huma.API, v any) *huma.Schema {
	return api.OpenAPI().Components.Schemas.Schema(
		reflect.TypeOf(v), true, "",
	)
}

// fieldSchema returns a huma.Schema based on a reflect.Type.
func fieldSchema(t reflect.Type) *huma.Schema {
	switch t.Kind() {
	case reflect.Bool:
		return &huma.Schema{Type: "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &huma.Schema{Type: "integer"}
	case reflect.Float32, reflect.Float64:
		return &huma.Schema{Type: "number"}
	default:
		return &huma.Schema{Type: "string"}
	}
}

// splitComma splits a string by comma.
func splitComma(s string) []string {
	var result []string
	start := 0
	for i, c := range s {
		if c == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	return append(result, s[start:])
}
