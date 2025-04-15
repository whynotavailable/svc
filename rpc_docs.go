package svc

import (
	"net/http"
	"reflect"
)

// This will be used to generate and deal with function metadata.

func (c *RpcContainer) GenerateDocs() {
	docs := map[string]FunctionDoc{}

	for key, f := range c.functions {
		var bodyInfo any = nil

		if f.bodyType != nil {
			bodyInfo = GenerateSchema(f.bodyType)
		}

		docs[key] = FunctionDoc{
			Body: bodyInfo,
			Meta: f.meta,
		}
	}

	c.mux.HandleFunc("GET /_info", func(w http.ResponseWriter, r *http.Request) {
		WriteJson(w, docs)
	})
}

type FunctionDoc struct {
	Body any
	Meta map[string]string
}

type SchemaObject struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
}

type SchemaMap struct {
	Type                 string `json:"type"`
	AdditionalProperties any    `json:"additionalProperties"`
}

type SchemaArray struct {
	Type  string `json:"type"`
	Items any    `json:"items"`
}

type SchemaField struct {
	Type string `json:"type"`
}

func GenerateSchema(elemType reflect.Type) any {
	if elemType.Kind() == reflect.Pointer {
		elemType = elemType.Elem()
	}

	if elemType.Kind() == reflect.Struct {
		schema := SchemaObject{
			Type:       "object",
			Properties: map[string]any{},
		}

		for i := range elemType.NumField() {
			propType := elemType.Field(i)
			setName, ok := propType.Tag.Lookup("json")
			if !ok {
				setName = propType.Name
			}
			schema.Properties[setName] = GenerateSchema(propType.Type)
		}

		return schema
	} else if elemType.Kind() == reflect.Map {
		return SchemaMap{
			Type:                 "object",
			AdditionalProperties: GenerateSchema(elemType.Elem()),
		}
	} else if elemType.Kind() == reflect.Slice {
		return SchemaArray{
			Type:  "array",
			Items: GenerateSchema(elemType.Elem()),
		}
	} else {
		return SchemaField{
			Type: translateKind(elemType.Kind().String()),
		}
	}
}

var kindMapping map[string]string = map[string]string{
	"interface": "object",
}

func translateKind(kind string) string {
	if mapping, ok := kindMapping[kind]; ok {
		return mapping
	}

	return kind
}

func (container *RpcContainer) ServeInfo(w http.ResponseWriter) {
}
