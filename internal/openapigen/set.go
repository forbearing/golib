package openapigen

import (
	"fmt"
	"reflect"
	"strconv"

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

func setCreate[M types.Model](pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	gen := openapi3gen.NewGenerator()
	name := typ.Elem().Name()

	reqSchemaRef, err := gen.NewSchemaRefForValue(*new(M), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	setupExample(reqSchemaRef)

	pathItem.Post = &openapi3.Operation{
		OperationID: operationID(consts.Create, typ),
		Summary:     summary(consts.Create, typ),
		Description: description(consts.Create, typ),
		Tags:        tags(consts.Create, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Create %s", name),
				Required:    true,
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			resp := openapi3.NewResponses()

			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
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

			schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
			if err != nil {
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
	addDefaultHeaders(pathItem.Post)
}

func setDelete[M types.Model](pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	pathItem.Delete = &openapi3.Operation{
		OperationID: operationID(consts.Delete, typ),
		Summary:     summary(consts.Delete, typ),
		Description: description(consts.Delete, typ),
		Tags:        tags(consts.Delete, typ),
		Parameters:  idParameters,
		Responses: func() *openapi3.Responses {
			schemaRef204, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef204 = new(openapi3.SchemaRef)
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
	addDefaultHeaders(pathItem.Delete)
}

func setUpdate[M types.Model](pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	reqSchemaRef, err := gen.NewSchemaRefForValue(*new(M), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	setupExample(reqSchemaRef)

	pathItem.Put = &openapi3.Operation{
		OperationID: operationID(consts.Update, typ),
		Summary:     summary(consts.Update, typ),
		Description: description(consts.Update, typ),
		Tags:        tags(consts.Update, typ),
		Parameters:  idParameters,
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("The %s data to update", name),
				Required:    true,
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
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
	addDefaultHeaders(pathItem.Put)
}

func setUpdatePartial[M types.Model](pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	reqSchemaRef, err := gen.NewSchemaRefForValue(*new(M), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	setupExample(reqSchemaRef)

	pathItem.Patch = &openapi3.Operation{
		OperationID: operationID(consts.UpdatePartial, typ),
		Summary:     summary(consts.UpdatePartial, typ),
		Description: description(consts.UpdatePartial, typ),
		Tags:        []string{name},
		Parameters:  idParameters,
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Partial fields of %s to update", name),
				Required:    true,
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
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
	addDefaultHeaders(pathItem.Patch)
}

func setList[M types.Model](pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	schemaRef, err := gen.NewSchemaRefForValue(*new(M), nil)
	if err != nil {
		zap.S().Error(err)
		schemaRef = new(openapi3.SchemaRef)
	}
	doc.Components.Schemas[name] = schemaRef

	// TODO: 从 schema 字段中读取查询信息
	pathItem.Get = &openapi3.Operation{
		OperationID: operationID(consts.List, typ),
		Summary:     summary(consts.List, typ),
		Description: description(consts.List, typ),
		Tags:        tags(consts.List, typ),
		// TODO: query parameters get from schema tags
		Parameters: []*openapi3.ParameterRef{
			{
				Value: &openapi3.Parameter{
					Name:     "page",
					In:       "query",
					Required: false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:    &openapi3.Types{openapi3.TypeInteger},
							Default: 1,
						},
					},
					Description: "Page number",
				},
			},
			{
				Value: &openapi3.Parameter{
					Name:     "pageSize",
					In:       "query",
					Required: false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:    &openapi3.Types{openapi3.TypeInteger},
							Default: 10,
						},
					},
					Description: "Number of items per page",
				},
			},
			// 可扩展更多 query 参数，例如过滤字段、排序等
		},
		Responses: func() *openapi3.Responses {
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiListResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
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
	addDefaultHeaders(pathItem.Get)
}

func setGet[M types.Model](pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	schemaRef, err := gen.NewSchemaRefForValue(*new(M), nil)
	if err != nil {
		zap.S().Error(err)
		schemaRef = new(openapi3.SchemaRef)
	}

	doc.Components.Schemas[name] = schemaRef
	pathItem.Get = &openapi3.Operation{
		OperationID: operationID(consts.Get, typ),
		Summary:     summary(consts.Get, typ),
		Description: description(consts.Get, typ),
		Tags:        tags(consts.Get, typ),
		Parameters:  idParameters,
		Responses: func() *openapi3.Responses {
			schemeRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemeRef200 = new(openapi3.SchemaRef)
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
	addDefaultHeaders(pathItem.Get)
}

func setImport[M types.Model](pathItem *openapi3.PathItem) {
}

func setExport[M types.Model](pathItem *openapi3.PathItem) {
}

func setBatchCreate[M types.Model](pathItem *openapi3.PathItem) {
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
	reqSchemaRef, err := gen.NewSchemaRefForValue(*new(apiBatchRequest[M]), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	setupBatchExample(reqSchemaRef)

	pathItem.Post = &openapi3.Operation{
		OperationID: operationID(consts.BatchCreate, typ),
		Summary:     summary(consts.BatchCreate, typ),
		Description: description(consts.BatchCreate, typ),
		Tags:        tags(consts.BatchCreate, typ),
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
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil)
			if err != nil {
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
	addDefaultHeaders(pathItem.Post)
}

func setBatchDelete[M types.Model](pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	reqSchemaName := name + "BatchDeleteRequest"
	doc.Components.Schemas[reqSchemaName] = &openapi3.SchemaRef{
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
		OperationID: operationID(consts.BatchDelete, typ),
		Summary:     summary(consts.BatchDelete, typ),
		Description: description(consts.BatchDelete, typ),
		Tags:        tags(consts.BatchDelete, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required:    true,
				Description: fmt.Sprintf("IDs of %s to delete", name),
				Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
					Ref: "#/components/schemas/" + reqSchemaName,
				}),
			},
		},
		Responses: func() *openapi3.Responses {
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
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
	addDefaultHeaders(pathItem.Delete)
}

func setBatchUpdate[M types.Model](pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	reqSchemaRef, err := gen.NewSchemaRefForValue(*new(apiBatchRequest[M]), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	setupBatchExample(reqSchemaRef)

	pathItem.Put = &openapi3.Operation{
		OperationID: operationID(consts.BatchUpdate, typ),
		Summary:     summary(consts.BatchUpdate, typ),
		Description: description(consts.BatchUpdate, typ),
		Tags:        tags(consts.BatchUpdate, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Request body for batch updating %s", name),
				Required:    true,
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
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
	addDefaultHeaders(pathItem.Put)
}

func setBatchUpdatePartial[M types.Model](pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()

	reqSchemaRef, err := gen.NewSchemaRefForValue(*new(apiBatchRequest[M]), nil)
	if err != nil {
		zap.S().Error(err)
		reqSchemaRef = new(openapi3.SchemaRef)
	}
	setupBatchExample(reqSchemaRef)

	pathItem.Patch = &openapi3.Operation{
		OperationID: operationID(consts.BatchUpdatePartial, typ),
		Summary:     summary(consts.BatchUpdatePartial, typ),
		Description: description(consts.BatchUpdatePartial, typ),
		Tags:        tags(consts.BatchUpdatePartial, typ),
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Description: fmt.Sprintf("Request body for batch partial updating %s", name),
				Required:    true,
				Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
			},
		},
		Responses: func() *openapi3.Responses {
			schemaRef200, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[M]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef200 = new(openapi3.SchemaRef)
			}
			schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil)
			if err != nil {
				zap.S().Error(err)
				schemaRef400 = new(openapi3.SchemaRef)
			}
			schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
			schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil)
			if err != nil {
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
	addDefaultHeaders(pathItem.Patch)
}

func addDefaultHeaders(op *openapi3.Operation) {
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

	// 避免重复添加
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

func operationID(op consts.HTTPVerb, typ reflect.Type) string {
	return fmt.Sprintf("%s%s", op, typ.Elem().Name())
}

func summary(op consts.HTTPVerb, typ reflect.Type) string {
	switch op {
	case consts.List, consts.BatchCreate, consts.BatchDelete, consts.BatchUpdate, consts.BatchUpdatePartial:
		return fmt.Sprintf("%s %s", op, pluralizeCli.Plural(typ.Elem().Name()))
	}
	return fmt.Sprintf("%s %s", op, typ.Elem().Name())
}

func description(op consts.HTTPVerb, typ reflect.Type) string {
	switch op {
	case consts.List, consts.BatchCreate, consts.BatchDelete, consts.BatchUpdate, consts.BatchUpdatePartial:
		return fmt.Sprintf("%s %s", op, pluralizeCli.Plural(typ.Elem().Name()))
	}
	return fmt.Sprintf("%s %s", op, typ.Elem().Name())
}

func tags(op consts.HTTPVerb, typ reflect.Type) []string {
	return []string{typ.Elem().Name()}
}

func exampleValue(code response.Code) map[string]any {
	return map[string]any{
		"code":       code.Code(),
		"data":       "null",
		"msg":        code.Msg(),
		"request_id": "string",
	}
}

// setupExample will remove field "created_at", "created_by", "updated_at", "updated_by"
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
//	  "id": "string",
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
		if k == "created_at" || k == "created_by" || k == "updated_at" || k == "updated_by" {
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
