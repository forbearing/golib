package openapigen

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/util"
	"github.com/gertd/go-pluralize"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"go.uber.org/zap"
)

var pluralizeCli = pluralize.NewClient()

var (
	success  = "success"
	idFormat = "" // eg: "uuid"
)

var (
	msgNotFound    = response.CodeNotFound.Msg()
	codeNotFound   = response.CodeNotFound.Code()
	statusNotFound = strconv.Itoa(response.CodeNotFound.Status())

	msgBadRequest    = response.CodeBadRequest.Msg()
	codeBadRequest   = response.CodeBadRequest.Code()
	statusBadRequest = strconv.Itoa(response.CodeBadRequest.Status())
)

var idParameters []*openapi3.ParameterRef = []*openapi3.ParameterRef{
	{
		Value: &openapi3.Parameter{
			In:       "path",
			Name:     "id",
			Required: true,
			Schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:   &openapi3.Types{openapi3.TypeString},
					Format: idFormat,
				},
			},
		},
	},
}

func setCreate[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	gen := openapi3gen.NewGenerator()
	name := typ.Elem().Name()

	var reqSchemaRef *openapi3.SchemaRef
	var err error
	if !model.IsModelEmpty[REQ]() {
		reqSchemaRef, err = gen.NewSchemaRefForValue(*new(REQ), nil)
		if err != nil {
			zap.S().Error(err)
			reqSchemaRef = new(openapi3.SchemaRef)
		}
	}
	setupExample(reqSchemaRef)
	addSchemaTitleDesc[M](reqSchemaRef)

	pathItem.Post = &openapi3.Operation{
		OperationID: operationID(consts.Create, typ),
		Summary:     summary(consts.Create, typ),
		Description: description(consts.Create, typ),
		Tags:        tags(path, consts.Create, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Create %s", name),
				Required:    !model.IsModelEmpty[REQ](),
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			resp := openapi3.NewResponses()
			var schemaRef200 *openapi3.SchemaRef
			var schemaRef400 *openapi3.SchemaRef
			var err error

			schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
					addSchemaTitleDesc[M](dataProperty)
				}
			}

			resp.Set("201", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s created", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
					// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
					// 	Ref: "#/components/schemas/" + typ.Elem().Name(),
					// }),
				},
			})

			if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf("Invalid request"),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			// resp.Set("401", &openapi3.ResponseRef{
			// 	Value: &openapi3.Response{
			// 		Description: util.ValueOf("Unauthorized"),
			// 		Content:     openapi3.NewContentWithJSONSchemaRef(errorSchemaRef),
			// 	},
			// })
			// resp.Set("409", &openapi3.ResponseRef{
			// 	Value: &openapi3.Response{
			// 		Description: util.ValueOf(fmt.Sprintf("%s already exists", name)),
			// 		Content:     openapi3.NewContentWithJSONSchemaRef(errorSchemaRef),
			// 	},
			// })
			// resp.Set("500", &openapi3.ResponseRef{
			// 	Value: &openapi3.Response{
			// 		Description: util.ValueOf("Internal server error"),
			// 		Content:     openapi3.NewContentWithJSONSchemaRef(errorSchemaRef),
			// 	},
			// })

			return resp
		}(),
	}
	addHeaderParameters(pathItem.Post)
}

func setDelete[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	pathItem.Delete = &openapi3.Operation{
		OperationID: operationID(consts.Delete, typ),
		Summary:     summary(consts.Delete, typ),
		Description: description(consts.Delete, typ),
		Tags:        tags(path, consts.Delete, typ),
		Parameters:  idParameters,
		Responses: func() *openapi3.Responses {
			schemaRef204, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef204 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef204.Value != nil && schemaRef204.Value.Properties != nil {
				if dataProperty, exists := schemaRef204.Value.Properties["data"]; exists {
					addSchemaTitleDesc[M](dataProperty)
				}
			}
			schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef204 = new(openapi3.SchemaRef)
			}
			schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
			resp := openapi3.NewResponses()
			resp.Set("204", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s deleted successfully", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef204),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf("Invalid request"),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})

			return resp
		}(),
	}
	addHeaderParameters(pathItem.Delete)
}

func setUpdate[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	var reqSchemaRef *openapi3.SchemaRef
	var err error
	if !model.IsModelEmpty[REQ]() {
		reqSchemaRef, err = gen.NewSchemaRefForValue(*new(REQ), nil)
		if err != nil {
			zap.S().Error(err)
			reqSchemaRef = new(openapi3.SchemaRef)
		}
	}
	setupExample(reqSchemaRef)
	addSchemaTitleDesc[M](reqSchemaRef)

	pathItem.Put = &openapi3.Operation{
		OperationID: operationID(consts.Update, typ),
		Summary:     summary(consts.Update, typ),
		Description: description(consts.Update, typ),
		Tags:        tags(path, consts.Update, typ),
		Parameters:  idParameters,
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("The %s data to update", name),
				Required:    !model.IsModelEmpty[REQ](),
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			var schemaRef200 *openapi3.SchemaRef
			var schemaRef400 *openapi3.SchemaRef
			var schemaRef404 *openapi3.SchemaRef
			var err error

			schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
					addSchemaTitleDesc[M](dataProperty)
				}
			}

			if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef404 = new(openapi3.SchemaRef)
			}
			schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
			resp := openapi3.NewResponses()
			resp.Set("200", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s updated successfully", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
					// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
					// 	Ref: "#/components/schemas/" + typ.Elem().Name(),
					// }),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf("Invalid request"),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})
			return resp
		}(),
	}
	addHeaderParameters(pathItem.Put)
}

func setUpdatePartial[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	var reqSchemaRef *openapi3.SchemaRef
	var err error
	if !model.IsModelEmpty[REQ]() {
		reqSchemaRef, err = gen.NewSchemaRefForValue(*new(REQ), nil)
		if err != nil {
			zap.S().Error(err)
			reqSchemaRef = new(openapi3.SchemaRef)
		}
	}
	setupExample(reqSchemaRef)
	addSchemaTitleDesc[M](reqSchemaRef)

	pathItem.Patch = &openapi3.Operation{
		OperationID: operationID(consts.Patch, typ),
		Summary:     summary(consts.Patch, typ),
		Description: description(consts.Patch, typ),
		Tags:        tags(path, consts.Patch, typ),
		Parameters:  idParameters,
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Partial fields of %s to update", name),
				Required:    !model.IsModelEmpty[REQ](),
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			var schemaRef200 *openapi3.SchemaRef
			var schemaRef400 *openapi3.SchemaRef
			var schemaRef404 *openapi3.SchemaRef

			schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
					addSchemaTitleDesc[M](dataProperty)
				}
			}
			if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef404 = new(openapi3.SchemaRef)
			}
			schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
			resp := openapi3.NewResponses()
			resp.Set("200", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s partially updated successfully", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
					// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
					// 	Ref: "#/components/schemas/" + typ.Elem().Name(),
					// }),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf("Invalid request"),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})
			return resp
		}(),
	}
	addHeaderParameters(pathItem.Patch)
}

func setList[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	schemaRef, err := gen.NewSchemaRefForValue(*new(M), nil)
	if err != nil {
		zap.S().Error(err)
		schemaRef = new(openapi3.SchemaRef)
	}
	// Add field descriptions to schema
	addSchemaTitleDesc[M](schemaRef)
	docMutex.Lock()
	doc.Components.Schemas[name] = schemaRef
	docMutex.Unlock()

	pathItem.Get = &openapi3.Operation{
		OperationID: operationID(consts.List, typ),
		Summary:     summary(consts.List, typ),
		Description: description(consts.List, typ),
		Tags:        tags(path, consts.List, typ),
		// Parameters: []*openapi3.ParameterRef{
		// 	{
		// 		Value: &openapi3.Parameter{
		// 			Name:     "page",
		// 			In:       "query",
		// 			Required: false,
		// 			Schema: &openapi3.SchemaRef{
		// 				Value: &openapi3.Schema{
		// 					Type:    &openapi3.Types{openapi3.TypeInteger},
		// 					Default: 1,
		// 				},
		// 			},
		// 			Description: "Page number",
		// 		},
		// 	},
		// 	{
		// 		Value: &openapi3.Parameter{
		// 			Name:     "pageSize",
		// 			In:       "query",
		// 			Required: false,
		// 			Schema: &openapi3.SchemaRef{
		// 				Value: &openapi3.Schema{
		// 					Type:    &openapi3.Types{openapi3.TypeInteger},
		// 					Default: 10,
		// 				},
		// 			},
		// 			Description: "Number of items per page",
		// 		},
		// 	},
		// 	// Can extend more query parameters, such as filter fields, sorting, etc.
		// },
		Responses: func() *openapi3.Responses {
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiListResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
								addSchemaTitleDesc[M](itemsProperty.Value.Items)
							}
						}
					}
				}
			}
			schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiListResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiListResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef404 = new(openapi3.SchemaRef)
			}
			schemaRef404.Value.Example = exampleValue(response.CodeNotFound)

			resp := openapi3.NewResponses()
			resp.Set("200", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("List of %s", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(msgBadRequest),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(msgNotFound),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})

			return resp
		}(),
	}
	addQueryParameters[M](pathItem.Get)
	addHeaderParameters(pathItem.Get)
}

func setGet[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	schemaRef, err := gen.NewSchemaRefForValue(*new(M), nil)
	if err != nil {
		zap.S().Error(err)
		schemaRef = new(openapi3.SchemaRef)
	}

	// Add field descriptions to schema
	addSchemaTitleDesc[M](schemaRef)

	docMutex.Lock()
	doc.Components.Schemas[name] = schemaRef
	docMutex.Unlock()
	pathItem.Get = &openapi3.Operation{
		OperationID: operationID(consts.Get, typ),
		Summary:     summary(consts.Get, typ),
		Description: description(consts.Get, typ),
		Tags:        tags(path, consts.Get, typ),
		Parameters:  idParameters,
		Responses: func() *openapi3.Responses {
			schemeRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemeRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemeRef200.Value != nil && schemeRef200.Value.Properties != nil {
				if dataProperty, exists := schemeRef200.Value.Properties["data"]; exists {
					addSchemaTitleDesc[M](dataProperty)
				}
			}
			schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef404 = new(openapi3.SchemaRef)
			}
			schemaRef404.Value.Example = exampleValue(response.CodeNotFound)

			resp := openapi3.NewResponses()
			resp.Set("200", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s detail", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemeRef200),
					// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
					// 	Ref: "#/components/schemas/" + name,
					// }),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})
			return resp
		}(),
	}
	addHeaderParameters(pathItem.Get)
}

func setCreateMany[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	// schemaRef, err := gen.NewSchemaRefForValue(*new(M), nil)
	// if err != nil {
	// 	zap.S().Error(err)
	// 	schemaRef = new(openapi3.SchemaRef)
	// }
	// doc.Components.Schemas[name] = schemaRef

	// // 定义 BatchCreateRequest schema
	// reqSchemaName := name + "BatchRequest"
	// reqSchemaRef := &openapi3.SchemaRef{
	// 	Value: &openapi3.Schema{
	// 		Type:     &openapi3.Types{openapi3.TypeObject},
	// 		Required: []string{"items"},
	// 		Properties: map[string]*openapi3.SchemaRef{
	// 			"items": {
	// 				Value: &openapi3.Schema{
	// 					Type:  &openapi3.Types{openapi3.TypeArray},
	// 					Items: &openapi3.SchemaRef{Ref: "#/components/schemas/" + name},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// doc.Components.Schemas[reqSchemaName] = reqSchemaRef

	var reqSchemaRef *openapi3.SchemaRef
	var err error
	reqSchemaRef, err = gen.NewSchemaRefForValue(*new(apiBatchRequest[REQ]), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	// Add field descriptions to request body schema
	if reqSchemaRef.Value != nil && reqSchemaRef.Value.Properties != nil {
		if itemsProperty, exists := reqSchemaRef.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
			addSchemaTitleDesc[M](itemsProperty.Value.Items)
		}
	}
	setupBatchExample(reqSchemaRef)

	pathItem.Post = &openapi3.Operation{
		OperationID: operationID(consts.CreateMany, typ),
		Summary:     summary(consts.CreateMany, typ),
		Description: description(consts.CreateMany, typ),
		Tags:        tags(path, consts.CreateMany, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Request body for batch creating %s", name),
				Required:    true,
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
				// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				// 	Ref: "#/components/schemas/" + reqSchemaName,
				// }),
			},
		},
		Responses: func() *openapi3.Responses {
			var schemaRef200 *openapi3.SchemaRef
			var schemaRef400 *openapi3.SchemaRef
			var schemaRef404 *openapi3.SchemaRef
			var err error

			schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[RSP]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
								addSchemaTitleDesc[M](itemsProperty.Value.Items)
							}
						}
					}
				}
			}
			if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef404 = new(openapi3.SchemaRef)
			}
			schemaRef404.Value.Example = exampleValue(response.CodeNotFound)

			resp := openapi3.NewResponses()
			resp.Set("201", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s created", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})

			return resp
		}(),
	}
	addHeaderParameters(pathItem.Post)
}

func setDeleteMany[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	reqSchemaRef := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:     &openapi3.Types{openapi3.TypeObject},
			Required: []string{"ids"},
			Properties: map[string]*openapi3.SchemaRef{
				"ids": {
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeArray},
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:   &openapi3.Types{openapi3.TypeString},
								Format: idFormat,
							},
						},
					},
				},
			},
		},
	}

	pathItem.Delete = &openapi3.Operation{
		OperationID: operationID(consts.DeleteMany, typ),
		Summary:     summary(consts.DeleteMany, typ),
		Description: description(consts.DeleteMany, typ),
		Tags:        tags(path, consts.DeleteMany, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required:    true,
				Description: fmt.Sprintf("IDs of %s to delete", name),
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists && dataProperty.Value != nil && dataProperty.Value.Properties != nil {
					if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
						addSchemaTitleDesc[M](itemsProperty.Value.Items)
					}
				}
			}
			schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef404 = new(openapi3.SchemaRef)
			}

			resp := openapi3.NewResponses()
			resp.Set("200", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s deleted", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})
			return resp
		}(),
	}
	addHeaderParameters(pathItem.Delete)
}

func setUpdateMany[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	var reqSchemaRef *openapi3.SchemaRef
	var err error
	reqSchemaRef, err = gen.NewSchemaRefForValue(*new(apiBatchRequest[REQ]), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	// Add field descriptions to request body schema
	if reqSchemaRef.Value != nil && reqSchemaRef.Value.Properties != nil {
		if itemsProperty, exists := reqSchemaRef.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
			addSchemaTitleDesc[M](itemsProperty.Value.Items)
		}
	}
	setupBatchExample(reqSchemaRef)

	pathItem.Put = &openapi3.Operation{
		OperationID: operationID(consts.UpdateMany, typ),
		Summary:     summary(consts.UpdateMany, typ),
		Description: description(consts.UpdateMany, typ),
		Tags:        tags(path, consts.UpdateMany, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Request body for batch updating %s", name),
				Required:    true,
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			var schemaRef200 *openapi3.SchemaRef
			var schemaRef400 *openapi3.SchemaRef
			var schemaRef404 *openapi3.SchemaRef

			schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[RSP]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
								addSchemaTitleDesc[M](itemsProperty.Value.Items)
							}
						}
					}
				}
			}
			if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef404 = new(openapi3.SchemaRef)
			}
			schemaRef404.Value.Example = exampleValue(response.CodeNotFound)

			resp := openapi3.NewResponses()
			resp.Set("200", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s updated", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})

			return resp
		}(),
	}
	addHeaderParameters(pathItem.Put)
}

func setPatchMany[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	var reqSchemaRef *openapi3.SchemaRef
	var err error
	reqSchemaRef, err = gen.NewSchemaRefForValue(*new(apiBatchRequest[REQ]), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	// Add field descriptions to request body schema
	if reqSchemaRef.Value != nil && reqSchemaRef.Value.Properties != nil {
		if itemsProperty, exists := reqSchemaRef.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
			addSchemaTitleDesc[M](itemsProperty.Value.Items)
		}
	}
	setupBatchExample(reqSchemaRef)

	pathItem.Patch = &openapi3.Operation{
		OperationID: operationID(consts.PatchMany, typ),
		Summary:     summary(consts.PatchMany, typ),
		Description: description(consts.PatchMany, typ),
		Tags:        tags(path, consts.PatchMany, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Request body for batch partial updating %s", name),
				Required:    true,
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			var schemaRef200 *openapi3.SchemaRef
			var schemaRef400 *openapi3.SchemaRef
			var schemaRef404 *openapi3.SchemaRef
			var err error

			schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[RSP]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			// Add field descriptions to response data schema
			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
								addSchemaTitleDesc[M](itemsProperty.Value.Items)
							}
						}
					}
				}
			}
			if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil); err != nil {
				zap.S().Error(err)
				schemaRef404 = new(openapi3.SchemaRef)
			}
			schemaRef404.Value.Example = exampleValue(response.CodeNotFound)

			resp := openapi3.NewResponses()
			resp.Set("200", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s partially updated", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
				},
			})
			resp.Set("400", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
				},
			})
			resp.Set("404", &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
				},
			})

			return resp
		}(),
	}
	addHeaderParameters(pathItem.Patch)
}

func setImport[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
}

func setExport[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
}

func addHeaderParameters(op *openapi3.Operation) {
	headers := []*openapi3.ParameterRef{
		{
			Value: &openapi3.Parameter{
				In:          "header",
				Name:        "Authorization",
				Description: "Authentication token (e.g. Bearer <token>)",
				Required:    false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
		{
			Value: &openapi3.Parameter{
				In:          "header",
				Name:        "X-Request-ID",
				Description: "Optional request ID for tracing",
				Required:    false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
		{
			Value: &openapi3.Parameter{
				In:          "header",
				Name:        "X-Client-Version",
				Description: "Client version (e.g. v1.2.3)",
				Required:    false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
		{
			Value: &openapi3.Parameter{
				In:          "header",
				Name:        "Accept-Language",
				Description: "Preferred language (e.g. zh-CN, en-US)",
				Required:    false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
	}

	// Avoid duplicate additions
	existing := map[string]bool{}
	for _, p := range op.Parameters {
		if p.Value != nil {
			existing[p.Value.Name] = true
		}
	}

	for _, header := range headers {
		if header.Value != nil && !existing[header.Value.Name] {
			op.Parameters = append(op.Parameters, header)
		}
	}
}

var (
	// Cache field descriptions of model.Base to avoid frequent parsing
	baseModelDocsCache map[string]string
	baseModelDocsOnce  sync.Once
)

// getBaseModelDocs gets field descriptions of model.Base (with caching)
func getBaseModelDocs() map[string]string {
	baseModelDocsOnce.Do(func() {
		baseModel := &model.Base{}
		baseModelDocsCache = parseModelDocs(baseModel)
	})
	return baseModelDocsCache
}

// addSchemaTitleDesc adds field descriptions to schema properties
func addSchemaTitleDesc[M types.Model](schemaRef *openapi3.SchemaRef) {
	if schemaRef == nil || schemaRef.Value == nil || schemaRef.Value.Properties == nil {
		return
	}

	// Get model field descriptions
	modelInstance := *new(M)
	modelDocs := parseModelDocs(modelInstance)

	// Get field descriptions of model.Base (using cache)
	baseDocs := getBaseModelDocs()

	// Create a mapping from JSON property names to field descriptions
	propertyDescriptions := make(map[string]string)

	// Process model fields
	typ := reflect.TypeOf(*new(M)).Elem()
	for i := range typ.NumField() {
		field := typ.Field(i)
		jsonTag := getFieldTag(field, consts.TAG_JSON)
		if len(jsonTag) == 0 {
			continue
		}

		// Get field descriptions from model documentation
		if description, exists := modelDocs[field.Name]; exists && description != "" {
			propertyDescriptions[jsonTag] = description
			// Debug: log the mapping
			// fmt.Printf("DEBUG: Field %s -> JSON %s -> Description: %s\n", field.Name, jsonTag, description)
		}
	}

	// Process Base model fields
	baseType := reflect.TypeOf(*new(model.Base))
	for i := range baseType.NumField() {
		field := baseType.Field(i)
		jsonTag := getFieldTag(field, consts.TAG_JSON)
		if len(jsonTag) == 0 {
			continue
		}

		// Get field descriptions from Base model documentation
		if description, exists := baseDocs[field.Name]; exists && description != "" {
			propertyDescriptions[jsonTag] = description
		}
	}

	// Add descriptions to schema properties
	for propName, propRef := range schemaRef.Value.Properties {
		if propRef.Value == nil {
			continue
		}

		// Set description if available
		if description, exists := propertyDescriptions[propName]; exists && description != "" {
			// Create a copy of the schema to avoid shared reference issues
			if propRef.Value != nil {
				// Create a new schema that's a copy of the original
				newSchema := *propRef.Value
				newSchema.Description = description
				newSchema.Title = description
				// Create a new SchemaRef and update the Properties map
				schemaRef.Value.Properties[propName] = &openapi3.SchemaRef{Value: &newSchema}
			}
		}
	}
}

func addQueryParameters[M types.Model](op *openapi3.Operation) {
	queries := make([]*openapi3.ParameterRef, 0)

	// Get model field descriptions
	modelInstance := *new(M)
	modelDocs := parseModelDocs(modelInstance)

	typ := reflect.TypeOf(*new(M)).Elem()
	for i := range typ.NumField() {
		field := typ.Field(i)
		schemaTag := getFieldTag(field, consts.TAG_SCHEMA)
		if len(schemaTag) == 0 {
			continue
		}

		// Get field descriptions from model documentation
		description := modelDocs[field.Name]

		queries = append(queries, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        schemaTag,
				In:          "query",
				Required:    false,
				Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: fieldType2openapiType(field)}},
				Description: description,
			},
		})
	}

	// Get field descriptions of model.Base (using cache)
	baseDocs := getBaseModelDocs()

	baseType := reflect.TypeOf(*new(model.Base))
	for i := range baseType.NumField() {
		field := baseType.Field(i)
		schemaTag := getFieldTag(field, consts.TAG_SCHEMA)
		if len(schemaTag) == 0 {
			continue
		}

		// Get field descriptions from Base model documentation
		description := baseDocs[field.Name]

		queries = append(queries, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        schemaTag,
				In:          "query",
				Required:    false,
				Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: fieldType2openapiType(field)}},
				Description: description,
			},
		})
	}

	// queries := []*openapi3.ParameterRef{
	// 	{
	// 		Value: &openapi3.Parameter{
	// 			Name:     "page",
	// 			In:       "query",
	// 			Required: false,
	// 			Schema: &openapi3.SchemaRef{
	// 				Value: &openapi3.Schema{
	// 					Type:    &openapi3.Types{openapi3.TypeInteger},
	// 					Default: 1,
	// 				},
	// 			},
	// 			Description: "Page number",
	// 		},
	// 	},
	// 	{
	// 		Value: &openapi3.Parameter{
	// 			Name:     "size",
	// 			In:       "query",
	// 			Required: false,
	// 			Schema: &openapi3.SchemaRef{
	// 				Value: &openapi3.Schema{
	// 					Type:    &openapi3.Types{openapi3.TypeInteger},
	// 					Default: 10,
	// 				},
	// 			},
	// 			Description: "Number of items per page",
	// 		},
	// 	},
	// }

	// Avoid duplicate additions
	existing := map[string]bool{}
	for _, p := range op.Parameters {
		if p.Value != nil {
			existing[p.Value.Name] = true
		}
	}

	for _, query := range queries {
		if query.Value != nil && !existing[query.Value.Name] {
			op.Parameters = append(op.Parameters, query)
		}
	}
}

func operationID(op consts.HTTPVerb, typ reflect.Type) string {
	return fmt.Sprintf("%s%s", op, typ.Elem().Name())
}

func summary(op consts.HTTPVerb, typ reflect.Type) string {
	switch op {
	case consts.List, consts.CreateMany, consts.DeleteMany, consts.UpdateMany, consts.PatchMany:
		return fmt.Sprintf("%s %s", op, pluralizeCli.Plural(typ.Elem().Name()))
	}
	return fmt.Sprintf("%s %s", op, typ.Elem().Name())
}

func description(op consts.HTTPVerb, typ reflect.Type) string {
	switch op {
	case consts.List, consts.CreateMany, consts.DeleteMany, consts.UpdateMany, consts.PatchMany:
		return fmt.Sprintf("%s %s", op, pluralizeCli.Plural(typ.Elem().Name()))
	}
	return fmt.Sprintf("%s %s", op, typ.Elem().Name())
}

func tags(path string, op consts.HTTPVerb, typ reflect.Type) []string {
	// return []string{typ.Elem().Name()}
	tag := strings.TrimPrefix(path, `/api/`)
	tag = strings.TrimSuffix(tag, `/batch`)
	tag = strings.TrimSuffix(tag, `/{id}`)
	tag = strings.ReplaceAll(tag, "/", "_")
	return []string{tag}
}

func exampleValue(code response.Code) map[string]any {
	return map[string]any{
		"code":       code.Code(),
		"data":       "null",
		"msg":        code.Msg(),
		"request_id": "string",
	}
}

// setupExample will remove field "created_at", "created_by", "updated_at", "updated_by", "id".
//
// Before:
//
//	{
//	  "created_at": "2025-04-19T19:22:55.434Z",
//	  "created_by": "string",
//	  "desc": "string",
//	  "id": "string",
//	  "member_count": 0,
//	  "name": "string",
//	  "order": 0,
//	  "remark": "string",
//	  "updated_at": "2025-04-19T19:22:55.434Z",
//	  "updated_by": "string"
//	}
//
// After:
//
//	{
//	  "desc": "string",
//	  "member_count": 0,
//	  "name": "string",
//	  "order": 0,
//	  "remark": "string"
//	}
func setupExample(schemaRef *openapi3.SchemaRef) {
	if schemaRef == nil {
		return
	}
	if schemaRef.Value == nil {
		schemaRef.Value = new(openapi3.Schema)
	}
	props := schemaRef.Value.Properties
	examples := make(map[string]any)
	for k, v := range props {
		if k == "created_at" || k == "created_by" || k == "updated_at" || k == "updated_by" || k == "id" {
			continue
		}
		if v.Value == nil || v.Value.Type == nil {
			continue
		}
		if v.Value.Type.Is(openapi3.TypeString) {
			examples[k] = "string"
		}
		if v.Value.Type.Is(openapi3.TypeInteger) {
			examples[k] = 0
		}
		if v.Value.Type.Is(openapi3.TypeNumber) {
			examples[k] = 0.0
		}
		if v.Value.Type.Is(openapi3.TypeBoolean) {
			examples[k] = false
		}
		if v.Value.Type.Is(openapi3.TypeArray) {
			examples[k] = []any{}
		}
		if v.Value.Type.Is(openapi3.TypeObject) {
			examples[k] = map[string]any{}
		}
		if v.Value.Type.Is(openapi3.TypeNull) {
			examples[k] = nil
		}
		schemaRef.Value.Example = examples
	}
}

// setupBatchExample will remove field "created_at", "created_by", "updated_at", "updated_by"
//
// Before:
//
//	{
//	  "items": [
//	    {
//	      "created_at": "2025-04-19T19:22:25.166Z",
//	      "created_by": "string",
//	      "desc": "string",
//	      "id": "string",
//	      "member_count": 0,
//	      "name": "string",
//	      "order": 0,
//	      "remark": "string",
//	      "updated_at": "2025-04-19T19:22:25.166Z",
//	      "updated_by": "string"
//	    }
//	  ]
//	}
//
// After:
//
//	{
//	  "items": [
//	    {
//	      "desc": "string",
//	      "id": "string",
//	      "member_count": 0,
//	      "name": "string",
//	      "order": 0,
//	      "remark": "string"
//	    }
//	  ]
//	}

func fieldType2openapiType(field reflect.StructField) *openapi3.Types {
	typ := field.Type

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		return &openapi3.Types{openapi3.TypeString}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &openapi3.Types{openapi3.TypeInteger}
	case reflect.Float32, reflect.Float64:
		return &openapi3.Types{openapi3.TypeNumber}
	case reflect.Bool:
		return &openapi3.Types{openapi3.TypeBoolean}
	case reflect.Array:
		return &openapi3.Types{openapi3.TypeArray}
	case reflect.Struct:
		// fmt.Println("----- field name", field.Name, field.Type.Kind())
		return &openapi3.Types{openapi3.TypeObject}
	default:
		// fmt.Println("----- field name", field.Name, field.Type.Kind())
		return &openapi3.Types{openapi3.TypeNull}
	}
}

func setupBatchExample(schemaRef *openapi3.SchemaRef) {
	if schemaRef == nil {
		return
	}
	if schemaRef.Value == nil {
		schemaRef.Value = new(openapi3.Schema)
	}
	props := schemaRef.Value.Properties
	for k, v := range props {
		if k == "items" && v.Value.Type.Is(openapi3.TypeArray) {
			example := make(map[string]any)
			for k, v := range v.Value.Items.Value.Properties {
				if k == "created_at" || k == "created_by" || k == "updated_at" || k == "updated_by" {
					continue
				}
				if v.Value == nil || v.Value.Type == nil {
					continue
				}
				if v.Value.Type.Is(openapi3.TypeString) {
					example[k] = "string"
				}
				if v.Value.Type.Is(openapi3.TypeInteger) {
					example[k] = 0
				}
				if v.Value.Type.Is(openapi3.TypeNumber) {
					example[k] = 0.0
				}
				if v.Value.Type.Is(openapi3.TypeBoolean) {
					example[k] = false
				}
				if v.Value.Type.Is(openapi3.TypeArray) {
					example[k] = []any{}
				}
				if v.Value.Type.Is(openapi3.TypeObject) {
					example[k] = map[string]any{}
				}
				if v.Value.Type.Is(openapi3.TypeNull) {
					example[k] = nil
				}
			}
			v.Value.Items.Value.Example = example
		}
	}
}

type apiBatchRequest[T any] struct {
	Items []T `json:"items"`
}

type apiResponse[T any] struct {
	Code      int    `json:"code"`
	Data      T      `json:"data"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
}

// newApiResponseRefWithData generate a openapi3.SchemaRef with custom data.
func newApiResponseRefWithData(data any) *openapi3.SchemaRef {
	dataSchemaRef, err := openapi3gen.NewSchemaRefForValue(data, nil)
	if err != nil {
		zap.S().Error(err)
		dataSchemaRef = new(openapi3.SchemaRef)
	}
	schema := &openapi3.Schema{
		Type: &openapi3.Types{openapi3.TypeObject},
		Properties: map[string]*openapi3.SchemaRef{
			"code":       {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}}},
			"msg":        {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
			"request_id": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
			"data":       dataSchemaRef, // ✅ Use dynamically generated type
		},
	}
	return &openapi3.SchemaRef{Value: schema}
}

type apiListResponse[T any] struct {
	Code      int         `json:"code"`
	Data      listData[T] `json:"data"`
	Msg       string      `json:"msg"`
	RequestID string      `json:"request_id"`
}
type listData[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}

type apiBatchResponse[T any] struct {
	Code      int          `json:"code"`
	Data      batchData[T] `json:"data"`
	Msg       string       `json:"msg"`
	RequestID string       `json:"request_id"`
}
type listSummary struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}

type batchData[T any] struct {
	Items   []T            `json:"items"`
	Options map[string]any `json:"options"`
	Summary listSummary    `json:"summary"`
}

// parameters:
//   - name: limit
//     in: query
//     required: false
//     schema:
//       type: integer
//
//   - name: Authorization
//     in: header
//     required: true
//     schema:
//       type: string
//
//   - name: id
//     in: path
//     required: true
//     schema:
//       type: string
//
//   - name: session_id
//     in: cookie
//     required: false
//     schema:
//       type: string
