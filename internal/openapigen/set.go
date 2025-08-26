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

var removeFieldMap = map[string]bool{
	"id":         true,
	"created_at": true,
	"created_by": true,
	"updated_at": true,
	"updated_by": true,
	"deleted_at": true,
	"deleted_by": true,
}

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
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_CREATE)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_CREATE)
	reqSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(REQ), nil)
	rspSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	// gen := openapi3gen.NewGenerator()
	// var reqSchemaRef *openapi3.SchemaRef
	// var err error
	// if !model.IsModelEmpty[REQ]() {
	// 	if reqSchemaRef, err = gen.NewSchemaRefForValue(*new(REQ), nil); err == nil {
	// 		setupExample(reqSchemaRef)
	// 		addSchemaTitleDesc[M](reqSchemaRef)
	// 	}
	// }

	pathItem.Post = &openapi3.Operation{
		OperationID: operationID(consts.Create, typ),
		Summary:     summary(consts.Create, typ),
		Description: description(consts.Create, typ),
		Tags:        tags(path, consts.Create, typ),
		RequestBody: newRequestBody(reqKey),
		Responses:   newResponses(201, rspKey),
		// RequestBody: &openapi3.RequestBodyRef{Ref: "#/components/requestBodies/" + reqKey},
		// Responses:   openapi3.NewResponses(openapi3.WithStatus(201, &openapi3.ResponseRef{Ref: "#/components/responses/" + rspKey})),

		// Responses: func() *openapi3.Responses {
		// 	resp := openapi3.NewResponses()
		// 	// var schemaRef200 *openapi3.SchemaRef
		// 	// // var schemaRef400 *openapi3.SchemaRef
		// 	// var err error
		// 	//
		// 	// if schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 	// 	// Add field descriptions to response data schema
		// 	// 	if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
		// 	// 		if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
		// 	// 			addSchemaTitleDesc[RSP](dataProperty)
		// 	// 		}
		// 	// 	}
		// 	// }
		//
		// 	resp.Set("201", &openapi3.ResponseRef{
		// 		Ref: "#/components/responses/" + rspKey,
		// 		// Value: &openapi3.Response{
		// 		// 	Description: util.ValueOf(fmt.Sprintf("%s created", name)),
		// 		// 	Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
		// 		// },
		// 	})
		//
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef400 = new(openapi3.SchemaRef)
		// 	// }
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf("Invalid request"),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		//
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("401", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf("Unauthorized"),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(errorSchemaRef),
		// 	// 	},
		// 	// })
		// 	// resp.Set("409", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s already exists", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(errorSchemaRef),
		// 	// 	},
		// 	// })
		// 	// resp.Set("500", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf("Internal server error"),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(errorSchemaRef),
		// 	// 	},
		// 	// })
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Post)
	removeFieldsFromRequestBody(pathItem.Post)
}

func setDelete[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_DELETE)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_DELETE)
	rspSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
	registerSchema[M, REQ, RSP](reqKey, rspKey, nil, rspSchemaRef)

	pathItem.Delete = &openapi3.Operation{
		OperationID: operationID(consts.Delete, typ),
		Summary:     summary(consts.Delete, typ),
		Description: description(consts.Delete, typ),
		Tags:        tags(path, consts.Delete, typ),
		Parameters:  idParameters,
		Responses:   newResponses(204, rspKey),
		// Responses: func() *openapi3.Responses {
		// 	var schemaRef204 *openapi3.SchemaRef
		// 	var err error
		// 	if schemaRef204, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 		// Add field descriptions to response data schema
		// 		if schemaRef204.Value != nil && schemaRef204.Value.Properties != nil {
		// 			if dataProperty, exists := schemaRef204.Value.Properties["data"]; exists {
		// 				addSchemaTitleDesc[RSP](dataProperty)
		// 			}
		// 		}
		// 	}
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
		// 	// if err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef400 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
		// 	// schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
		// 	// if err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef204 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("204", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("%s deleted successfully", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef204),
		// 		},
		// 	})
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf("Invalid request"),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		//
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Delete)
}

func setUpdate[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_UPDATE)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_UPDATE)
	reqSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(REQ), nil)
	rspSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	var err error
	if !model.IsModelEmpty[REQ]() {
		if reqSchemaRef, err = gen.NewSchemaRefForValue(*new(REQ), nil); err == nil {
			setupExample(reqSchemaRef)
			addSchemaTitleDesc[M](reqSchemaRef)
		}
	}

	pathItem.Put = &openapi3.Operation{
		OperationID: operationID(consts.Update, typ),
		Summary:     summary(consts.Update, typ),
		Description: description(consts.Update, typ),
		Tags:        tags(path, consts.Update, typ),
		Parameters:  idParameters,
		RequestBody: newRequestBody(reqKey),
		Responses:   newResponses(200, rspKey),
		// RequestBody: &openapi3.RequestBodyRef{
		// 	Value: &openapi3.RequestBody{
		// 		Description: fmt.Sprintf("The %s data to update", name),
		// 		Required:    !model.IsModelEmpty[REQ](),
		// 		Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
		// 	},
		// },
		// Responses: func() *openapi3.Responses {
		// 	var schemaRef200 *openapi3.SchemaRef
		// 	// var schemaRef400 *openapi3.SchemaRef
		// 	// var schemaRef404 *openapi3.SchemaRef
		// 	var err error
		//
		// 	if schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 		// Add field descriptions to response data schema
		// 		if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
		// 			if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
		// 				addSchemaTitleDesc[RSP](dataProperty)
		// 			}
		// 		}
		// 	}
		//
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef400 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
		// 	// if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef404 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
		//
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("200", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("%s updated successfully", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
		// 			// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
		// 			// 	Ref: "#/components/schemas/" + typ.Elem().Name(),
		// 			// }),
		// 		},
		// 	})
		//
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf("Invalid request"),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Put)
	removeFieldsFromRequestBody(pathItem.Put)
}

func setPatch[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_PATCH)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_PATCH)
	reqSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(REQ), nil)
	rspSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	var err error
	if !model.IsModelEmpty[REQ]() {
		if reqSchemaRef, err = gen.NewSchemaRefForValue(*new(REQ), nil); err == nil {
			setupExample(reqSchemaRef)
			addSchemaTitleDesc[M](reqSchemaRef)
		}
	}

	pathItem.Patch = &openapi3.Operation{
		OperationID: operationID(consts.Patch, typ),
		Summary:     summary(consts.Patch, typ),
		Description: description(consts.Patch, typ),
		Tags:        tags(path, consts.Patch, typ),
		Parameters:  idParameters,
		RequestBody: newRequestBody(reqKey),
		Responses:   newResponses(200, rspKey),
		// RequestBody: &openapi3.RequestBodyRef{
		// 	Value: &openapi3.RequestBody{
		// 		Description: fmt.Sprintf("Partial fields of %s to update", name),
		// 		Required:    !model.IsModelEmpty[REQ](),
		// 		Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
		// 	},
		// },
		// Responses: func() *openapi3.Responses {
		// 	var schemaRef200 *openapi3.SchemaRef
		// 	// var schemaRef400 *openapi3.SchemaRef
		// 	// var schemaRef404 *openapi3.SchemaRef
		//
		// 	if schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 		// Add field descriptions to response data schema
		// 		if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
		// 			if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
		// 				addSchemaTitleDesc[RSP](dataProperty)
		// 			}
		// 		}
		// 	}
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef400 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
		// 	// if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef404 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("200", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("%s partially updated successfully", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
		// 			// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
		// 			// 	Ref: "#/components/schemas/" + typ.Elem().Name(),
		// 			// }),
		// 		},
		// 	})
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf("Invalid request"),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Patch)
}

func setList[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_LIST)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_LIST)
	reqSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(REQ), nil)

	var rspSchemaRef *openapi3.SchemaRef
	var err error
	if model.AreTypesEqual[M, REQ, RSP]() {
		if rspSchemaRef, err = openapi3gen.NewSchemaRefForValue(*new(apiListResponse[M]), nil); err == nil {
			// Add field descriptions to response data schema
			if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
				if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
								addSchemaTitleDesc[M](itemsProperty.Value.Items)
							}
						}
					}
				}
			}
		}
	} else {
		if rspSchemaRef, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
			if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
				if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
					addSchemaTitleDesc[RSP](dataProperty)
				}
			}
		}
	}
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	pathItem.Get = &openapi3.Operation{
		OperationID: operationID(consts.List, typ),
		Summary:     summary(consts.List, typ),
		Description: description(consts.List, typ),
		Tags:        tags(path, consts.List, typ),
		Responses:   newResponses(200, rspKey),
		// // Parameters: []*openapi3.ParameterRef{
		// // 	{
		// // 		Value: &openapi3.Parameter{
		// // 			Name:     "page",
		// // 			In:       "query",
		// // 			Required: false,
		// // 			Schema: &openapi3.SchemaRef{
		// // 				Value: &openapi3.Schema{
		// // 					Type:    &openapi3.Types{openapi3.TypeInteger},
		// // 					Default: 1,
		// // 				},
		// // 			},
		// // 			Description: "Page number",
		// // 		},
		// // 	},
		// // 	{
		// // 		Value: &openapi3.Parameter{
		// // 			Name:     "pageSize",
		// // 			In:       "query",
		// // 			Required: false,
		// // 			Schema: &openapi3.SchemaRef{
		// // 				Value: &openapi3.Schema{
		// // 					Type:    &openapi3.Types{openapi3.TypeInteger},
		// // 					Default: 10,
		// // 				},
		// // 			},
		// // 			Description: "Number of items per page",
		// // 		},
		// // 	},
		// // 	// Can extend more query parameters, such as filter fields, sorting, etc.
		// // },
		// Responses: func() *openapi3.Responses {
		// 	var schemaRef200 *openapi3.SchemaRef
		// 	var err error
		// 	if model.AreTypesEqual[M, REQ, RSP]() {
		// 		if schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiListResponse[M]), nil); err == nil {
		// 			// Add field descriptions to response data schema
		// 			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
		// 					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
		// 						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
		// 							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
		// 								addSchemaTitleDesc[M](itemsProperty.Value.Items)
		// 							}
		// 						}
		// 					}
		// 				}
		// 			}
		// 		}
		// 	} else {
		// 		if !model.IsModelEmpty[RSP]() {
		// 			if schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 				if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
		// 					if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
		// 						addSchemaTitleDesc[RSP](dataProperty)
		// 					}
		// 				}
		// 			}
		// 		}
		// 	}
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiListResponse[string]), nil)
		// 	// if err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef400 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
		// 	// schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiListResponse[string]), nil)
		// 	// if err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef404 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
		//
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("200", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("List of %s", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
		// 		},
		// 	})
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(msgBadRequest),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(msgNotFound),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		//
		// 	return resp
		// }(),
	}
	addQueryParameters[M, REQ, RSP](pathItem.Get)
	addHeaderParameters(pathItem.Get)
}

func setGet[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_GET)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_GET)
	reqSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(REQ), nil)
	rspSchemaRef, _ := openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	var schemaRef *openapi3.SchemaRef
	var err error
	if schemaRef, err = gen.NewSchemaRefForValue(*new(M), nil); err == nil {
		// Add field descriptions to schema
		addSchemaTitleDesc[M](schemaRef)
	}

	pathItem.Get = &openapi3.Operation{
		OperationID: operationID(consts.Get, typ),
		Summary:     summary(consts.Get, typ),
		Description: description(consts.Get, typ),
		Tags:        tags(path, consts.Get, typ),
		Parameters:  idParameters,
		Responses:   newResponses(200, rspKey),
		// Responses: func() *openapi3.Responses {
		// 	var schemaRef200 *openapi3.SchemaRef
		// 	var err error
		// 	if schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 		// Add field descriptions to response data schema
		// 		if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
		// 			if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
		// 				addSchemaTitleDesc[RSP](dataProperty)
		// 			}
		// 		}
		// 	}
		//
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
		// 	// if err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef400 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
		// 	// schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
		// 	// if err != nil {
		// 	// 	zap.S().Error(err)
		// 	// 	schemaRef404 = new(openapi3.SchemaRef)
		// 	// }
		// 	// schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
		//
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("200", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("%s detail", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
		// 			// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
		// 			// 	Ref: "#/components/schemas/" + name,
		// 			// }),
		// 		},
		// 	})
		//
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Get)
}

func setCreateMany[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_CREATE_MANY)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_CREATE_MANY)

	var reqSchemaRef *openapi3.SchemaRef
	var rspSchemaRef *openapi3.SchemaRef
	if model.AreTypesEqual[M, REQ, RSP]() {
		reqSchemaRef, _ = openapi3gen.NewSchemaRefForValue(*new(apiBatchRequest[REQ]), nil)
		if reqSchemaRef.Value != nil && reqSchemaRef.Value.Properties != nil {
			if itemsProperty, exists := reqSchemaRef.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
				addSchemaTitleDesc[M](itemsProperty.Value.Items)
			}
		}
		rspSchemaRef, _ = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[RSP]), nil)
		if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
			if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
				if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
					if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
						if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
							addSchemaTitleDesc[RSP](itemsProperty.Value.Items)
						}
					}
				}
			}
		}
	} else {
		reqSchemaRef, _ = openapi3gen.NewSchemaRefForValue(*new(REQ), nil)
		rspSchemaRef, _ = openapi3gen.NewSchemaRefForValue(*new(RSP), nil)
		if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
			if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
				addSchemaTitleDesc[RSP](dataProperty)
			}
		}
	}
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	// // // 定义 BatchCreateRequest schema
	// // reqSchemaName := name + "BatchRequest"
	// // reqSchemaRef := &openapi3.SchemaRef{
	// // 	Value: &openapi3.Schema{
	// // 		Type:     &openapi3.Types{openapi3.TypeObject},
	// // 		Required: []string{"items"},
	// // 		Properties: map[string]*openapi3.SchemaRef{
	// // 			"items": {
	// // 				Value: &openapi3.Schema{
	// // 					Type:  &openapi3.Types{openapi3.TypeArray},
	// // 					Items: &openapi3.SchemaRef{Ref: "#/components/schemas/" + name},
	// // 				},
	// // 			},
	// // 		},
	// // 	},
	// // }
	// // doc.Components.Schemas[reqSchemaName] = reqSchemaRef
	//
	// var err error
	// var reqSchemaRef *openapi3.SchemaRef
	// if reqSchemaRef, err = gen.NewSchemaRefForValue(*new(apiBatchRequest[REQ]), nil); err == nil {
	// 	// Add field descriptions to request body schema
	// 	if reqSchemaRef.Value != nil && reqSchemaRef.Value.Properties != nil {
	// 		if itemsProperty, exists := reqSchemaRef.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
	// 			addSchemaTitleDesc[M](itemsProperty.Value.Items)
	// 		}
	// 	}
	// 	setupBatchExample(reqSchemaRef)
	// }

	pathItem.Post = &openapi3.Operation{
		OperationID: operationID(consts.CreateMany, typ),
		Summary:     summary(consts.CreateMany, typ),
		Description: description(consts.CreateMany, typ),
		Tags:        tags(path, consts.CreateMany, typ),
		RequestBody: newRequestBody(reqKey),
		Responses:   newResponses(201, rspKey),
		// RequestBody: &openapi3.RequestBodyRef{
		// 	Value: &openapi3.RequestBody{
		// 		Description: fmt.Sprintf("Request body for batch creating %s", name),
		// 		Required:    true,
		// 		Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
		// 		// Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
		// 		// 	Ref: "#/components/schemas/" + reqSchemaName,
		// 		// }),
		// 	},
		// },
		// Responses: func() *openapi3.Responses {
		// 	var rspSchemaRef200 *openapi3.SchemaRef
		// 	// var schemaRef400 *openapi3.SchemaRef
		// 	// var schemaRef404 *openapi3.SchemaRef
		// 	var err error
		//
		// 	if model.AreTypesEqual[M, REQ, RSP]() {
		// 		if rspSchemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[M]), nil); err == nil {
		// 			// Add field descriptions to response data schema
		// 			if rspSchemaRef200.Value != nil && rspSchemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := rspSchemaRef200.Value.Properties["data"]; exists {
		// 					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
		// 						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
		// 							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
		// 								addSchemaTitleDesc[M](itemsProperty.Value.Items)
		// 							}
		// 						}
		// 					}
		// 				}
		// 			}
		// 		}
		// 		// // Mybe used in the future, DO NOT DELETE it.
		// 		// if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil); err != nil {
		// 		// 	zap.S().Error(err)
		// 		// 	schemaRef400 = new(openapi3.SchemaRef)
		// 		// }
		// 		// schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
		// 		// if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil); err != nil {
		// 		// 	zap.S().Error(err)
		// 		// 	schemaRef404 = new(openapi3.SchemaRef)
		// 		// }
		// 		// schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
		// 	} else {
		// 		if rspSchemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 			if rspSchemaRef200.Value != nil && rspSchemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := rspSchemaRef200.Value.Properties["data"]; exists {
		// 					addSchemaTitleDesc[RSP](dataProperty)
		// 				}
		// 			}
		// 		}
		// 	}
		//
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("201", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("%s created", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(rspSchemaRef200),
		// 		},
		// 	})
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		//
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Post)
	removeFieldsFromBatchRequestBody(pathItem.Post)
}

func setDeleteMany[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.DeleteMany)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.DeleteMany)
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
	var rspSchemaRef *openapi3.SchemaRef
	if model.AreTypesEqual[M, REQ, RSP]() {
		rspSchemaRef, _ = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[RSP]), nil)
		if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
			if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists && dataProperty.Value != nil && dataProperty.Value.Properties != nil {
				if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
					addSchemaTitleDesc[RSP](itemsProperty.Value.Items)
				}
			}
		}
	} else {
		rspSchemaRef, _ = openapi3gen.NewSchemaRefForValue(*new(RSP), nil)
		if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
			if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
				addSchemaTitleDesc[RSP](dataProperty)
			}
		}
	}
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	pathItem.Delete = &openapi3.Operation{
		OperationID: operationID(consts.DeleteMany, typ),
		Summary:     summary(consts.DeleteMany, typ),
		Description: description(consts.DeleteMany, typ),
		Tags:        tags(path, consts.DeleteMany, typ),
		RequestBody: newRequestBody(reqKey),
		Responses:   newResponses(204, rspKey),
		// RequestBody: &openapi3.RequestBodyRef{
		// 	Value: &openapi3.RequestBody{
		// 		Required:    true,
		// 		Description: fmt.Sprintf("IDs of %s to delete", name),
		// 		Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
		// 	},
		// },
		// Responses: func() *openapi3.Responses {
		// 	var schemaRef200 *openapi3.SchemaRef
		// 	var err error
		//
		// 	if model.AreTypesEqual[M, REQ, RSP]() {
		// 		if schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[M]), nil); err == nil {
		// 			// Add field descriptions to response data schema
		// 			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists && dataProperty.Value != nil && dataProperty.Value.Properties != nil {
		// 					if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
		// 						addSchemaTitleDesc[M](itemsProperty.Value.Items)
		// 					}
		// 				}
		// 			}
		// 		}
		// 		// // Mybe used in the future, DO NOT DELETE it.
		// 		// schemaRef400, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
		// 		// schemaRef404, err := openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil)
		// 	} else {
		// 		if schemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 			if schemaRef200.Value != nil && schemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := schemaRef200.Value.Properties["data"]; exists {
		// 					addSchemaTitleDesc[RSP](dataProperty)
		// 				}
		// 			}
		// 		}
		// 	}
		//
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("200", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("%s deleted", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef200),
		// 		},
		// 	})
		//
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Delete)
}

func setUpdateMany[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_UPDATE_MANY)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_UPDATE_MANY)

	var reqSchemaRef *openapi3.SchemaRef
	var rspSchemaRef *openapi3.SchemaRef
	if model.AreTypesEqual[M, REQ, RSP]() {
		reqSchemaRef, _ = gen.NewSchemaRefForValue(*new(apiBatchRequest[REQ]), nil)
		if reqSchemaRef.Value != nil && reqSchemaRef.Value.Properties != nil {
			if itemsProperty, exists := reqSchemaRef.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
				addSchemaTitleDesc[M](itemsProperty.Value.Items)
			}
		}
		rspSchemaRef, _ = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[REQ]), nil)
		if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
			if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
				if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
					if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
						if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
							addSchemaTitleDesc[REQ](itemsProperty.Value.Items)
						}
					}
				}
			}
		}
	} else {
		rspSchemaRef, _ = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil)
		if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
			if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
				addSchemaTitleDesc[RSP](dataProperty)
			}
		}
	}
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	pathItem.Put = &openapi3.Operation{
		OperationID: operationID(consts.UpdateMany, typ),
		Summary:     summary(consts.UpdateMany, typ),
		Description: description(consts.UpdateMany, typ),
		Tags:        tags(path, consts.UpdateMany, typ),
		RequestBody: newRequestBody(reqKey),
		Responses:   newResponses(200, rspKey),
		// RequestBody: &openapi3.RequestBodyRef{
		// 	Value: &openapi3.RequestBody{
		// 		Description: fmt.Sprintf("Request body for batch updating %s", name),
		// 		Required:    true,
		// 		Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
		// 	},
		// },
		// Responses: func() *openapi3.Responses {
		// 	var rspSchemaRef200 *openapi3.SchemaRef
		// 	// var schemaRef400 *openapi3.SchemaRef
		// 	// var schemaRef404 *openapi3.SchemaRef
		//
		// 	if model.AreTypesEqual[M, REQ, RSP]() {
		// 		if rspSchemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[RSP]), nil); err == nil {
		// 			// Add field descriptions to response data schema
		// 			if rspSchemaRef200.Value != nil && rspSchemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := rspSchemaRef200.Value.Properties["data"]; exists {
		// 					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
		// 						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
		// 							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
		// 								addSchemaTitleDesc[M](itemsProperty.Value.Items)
		// 							}
		// 						}
		// 					}
		// 				}
		// 			}
		// 		}
		// 		// // Mybe used in the future, DO NOT DELETE it.
		// 		// if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
		// 		// 	zap.S().Error(err)
		// 		// 	schemaRef400 = new(openapi3.SchemaRef)
		// 		// }
		// 		// schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
		// 		// if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err != nil {
		// 		// 	zap.S().Error(err)
		// 		// 	schemaRef404 = new(openapi3.SchemaRef)
		// 		// }
		// 		// schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
		// 	} else {
		// 		if rspSchemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[RSP]), nil); err == nil {
		// 			// Add field descriptions to response data schema
		// 			if rspSchemaRef200.Value != nil && rspSchemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := rspSchemaRef200.Value.Properties["data"]; exists {
		// 					addSchemaTitleDesc[RSP](dataProperty)
		// 				}
		// 			}
		// 		}
		// 	}
		// 	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef200)
		//
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("200", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("%s updated", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(rspSchemaRef200),
		// 		},
		// 	})
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		//
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Put)
	removeFieldsFromBatchRequestBody(pathItem.Put)
}

func setPatchMany[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	gen := openapi3gen.NewGenerator()
	typ := reflect.TypeOf(*new(M))
	name := typ.Elem().Name()
	reqKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_PATCH_MANY)
	rspKey := fmt.Sprintf("%s_%s", strings.ToLower(name), consts.PHASE_PATCH_MANY)

	var reqSchemaRef *openapi3.SchemaRef
	var rspSchemaRef *openapi3.SchemaRef
	if model.AreTypesEqual[M, REQ, RSP]() {
		reqSchemaRef, _ = gen.NewSchemaRefForValue(*new(apiBatchRequest[REQ]), nil)
		if reqSchemaRef.Value != nil && reqSchemaRef.Value.Properties != nil {
			if itemsProperty, exists := reqSchemaRef.Value.Properties["items"]; exists && itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
				addSchemaTitleDesc[M](itemsProperty.Value.Items)
			}
		}
		rspSchemaRef, _ = gen.NewSchemaRefForValue(*new(apiBatchResponse[RSP]), nil)
		if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
			if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
				if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
					if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
						if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
							addSchemaTitleDesc[M](itemsProperty.Value.Items)
						}
					}
				}
			}
		}
	} else {
		reqSchemaRef, _ = gen.NewSchemaRefForValue(*new(REQ), nil)
		rspSchemaRef, _ = gen.NewSchemaRefForValue(*new(RSP), nil)
		if rspSchemaRef.Value != nil && rspSchemaRef.Value.Properties != nil {
			if dataProperty, exists := rspSchemaRef.Value.Properties["data"]; exists {
				addSchemaTitleDesc[RSP](dataProperty)
			}
		}
	}
	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef)

	pathItem.Patch = &openapi3.Operation{
		OperationID: operationID(consts.PatchMany, typ),
		Summary:     summary(consts.PatchMany, typ),
		Description: description(consts.PatchMany, typ),
		Tags:        tags(path, consts.PatchMany, typ),
		RequestBody: newRequestBody(reqKey),
		Responses:   newResponses(200, rspKey),
		// RequestBody: &openapi3.RequestBodyRef{
		// 	Value: &openapi3.RequestBody{
		// 		Description: fmt.Sprintf("Request body for batch partial updating %s", name),
		// 		Required:    true,
		// 		Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
		// 	},
		// },
		// Responses: func() *openapi3.Responses {
		// 	var rspSchemaRef200 *openapi3.SchemaRef
		// 	// var schemaRef400 *openapi3.SchemaRef
		// 	// var schemaRef404 *openapi3.SchemaRef
		// 	var err error
		//
		// 	if model.AreTypesEqual[M, REQ, RSP]() {
		// 		if rspSchemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[RSP]), nil); err == nil {
		// 			// Add field descriptions to response data schema
		// 			if rspSchemaRef200.Value != nil && rspSchemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := rspSchemaRef200.Value.Properties["data"]; exists {
		// 					if dataProperty.Value != nil && dataProperty.Value.Properties != nil {
		// 						if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
		// 							if itemsProperty.Value != nil && itemsProperty.Value.Items != nil {
		// 								addSchemaTitleDesc[M](itemsProperty.Value.Items)
		// 							}
		// 						}
		// 					}
		// 				}
		// 			}
		// 		}
		// 		// // Mybe used in the future, DO NOT DELETE it.
		// 		// if schemaRef400, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil); err != nil {
		// 		// 	zap.S().Error(err)
		// 		// 	schemaRef400 = new(openapi3.SchemaRef)
		// 		// }
		// 		// schemaRef400.Value.Example = exampleValue(response.CodeBadRequest)
		// 		// if schemaRef404, err = openapi3gen.NewSchemaRefForValue(*new(apiBatchResponse[string]), nil); err != nil {
		// 		// 	zap.S().Error(err)
		// 		// 	schemaRef404 = new(openapi3.SchemaRef)
		// 		// }
		// 		// schemaRef404.Value.Example = exampleValue(response.CodeNotFound)
		// 	} else {
		// 		if rspSchemaRef200, err = openapi3gen.NewSchemaRefForValue(*new(apiResponse[string]), nil); err == nil {
		// 			if rspSchemaRef200.Value != nil && rspSchemaRef200.Value.Properties != nil {
		// 				if dataProperty, exists := rspSchemaRef200.Value.Properties["data"]; exists {
		// 					addSchemaTitleDesc[RSP](dataProperty)
		// 				}
		// 			}
		// 		}
		// 	}
		//
		// 	registerSchema[M, REQ, RSP](reqKey, rspKey, reqSchemaRef, rspSchemaRef200)
		// 	resp := openapi3.NewResponses()
		// 	resp.Set("200", &openapi3.ResponseRef{
		// 		Value: &openapi3.Response{
		// 			Description: util.ValueOf(fmt.Sprintf("%s partially updated", name)),
		// 			Content:     openapi3.NewContentWithJSONSchemaRef(rspSchemaRef200),
		// 		},
		// 	})
		// 	// // Mybe used in the future, DO NOT DELETE it.
		// 	// resp.Set("400", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef400),
		// 	// 	},
		// 	// })
		// 	// resp.Set("404", &openapi3.ResponseRef{
		// 	// 	Value: &openapi3.Response{
		// 	// 		Description: util.ValueOf(fmt.Sprintf("%s not found", name)),
		// 	// 		Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef404),
		// 	// 	},
		// 	// })
		//
		// 	return resp
		// }(),
	}
	addHeaderParameters(pathItem.Patch)
	removeFieldsFromBatchRequestBody(pathItem.Patch)
}

func setImport[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	// pathItem.Post = &openapi3.Operation{
	// 	OperationID: "import" + reflect.TypeOf(*new(M)).Elem().Name(),
	// 	Summary:     "Import " + reflect.TypeOf(*new(M)).Elem().Name() + " data",
	// 	Description: "Import data from CSV/Excel file",
	// 	Tags:        tags(path, "import", reflect.TypeOf(*new(M))),
	// 	RequestBody: &openapi3.RequestBodyRef{
	// 		Value: &openapi3.RequestBody{
	// 			Description: "File to import",
	// 			Required:    true,
	// 			Content: openapi3.Content{
	// 				"multipart/form-data": &openapi3.MediaType{
	// 					Schema: &openapi3.SchemaRef{
	// 						Value: &openapi3.Schema{
	// 							Type: &openapi3.Types{openapi3.TypeObject},
	// 							Properties: map[string]*openapi3.SchemaRef{
	// 								"file": {
	// 									Value: &openapi3.Schema{
	// 										Type:   &openapi3.Types{openapi3.TypeString},
	// 										Format: "binary",
	// 									},
	// 								},
	// 							},
	// 							Required: []string{"file"},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Responses: newResponses(200, "ImportResponse"),
	// }
}

func setExport[M types.Model, REQ types.Request, RSP types.Response](path string, pathItem *openapi3.PathItem) {
	// pathItem.Get = &openapi3.Operation{
	// 	OperationID: "export" + reflect.TypeOf(*new(M)).Elem().Name(),
	// 	Summary:     "Export " + reflect.TypeOf(*new(M)).Elem().Name() + " data",
	// 	Description: "Export data to CSV/Excel file",
	// 	Tags:        tags(path, "export", reflect.TypeOf(*new(M))),
	// 	Parameters: []*openapi3.ParameterRef{
	// 		{
	// 			Value: &openapi3.Parameter{
	// 				Name:        "format",
	// 				In:          "query",
	// 				Description: "Export format",
	// 				Schema: &openapi3.SchemaRef{
	// 					Value: &openapi3.Schema{
	// 						Type:    &openapi3.Types{openapi3.TypeString},
	// 						Enum:    []any{"csv", "xlsx"},
	// 						Default: "csv",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Responses: &openapi3.Responses{
	// 		MapOfResponseOrRefValues: openapi3.ResponsesMap{
	// 			"200": &openapi3.ResponseRef{
	// 				Value: &openapi3.Response{
	// 					Description: util.ValueOf("Export file"),
	// 					Content: openapi3.Content{
	// 						"text/csv": &openapi3.MediaType{
	// 							Schema: &openapi3.SchemaRef{
	// 								Value: &openapi3.Schema{
	// 									Type:   &openapi3.Types{openapi3.TypeString},
	// 									Format: "binary",
	// 								},
	// 							},
	// 						},
	// 						"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": &openapi3.MediaType{
	// 							Schema: &openapi3.SchemaRef{
	// 								Value: &openapi3.Schema{
	// 									Type:   &openapi3.Types{openapi3.TypeString},
	// 									Format: "binary",
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
}

// register Model, Model Payload, Model Result into openapi3 schema.
func registerSchema[M types.Model, REQ types.Request, RSP types.Response](reqKey, rspKey string, reqSchemaRef *openapi3.SchemaRef, rspSchemaRef *openapi3.SchemaRef) {
	if !model.IsModelEmpty[M]() {
		typ := reflect.TypeOf(*new(M))
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		name := typ.Name()
		docMutex.Lock()
		if doc.Components.Schemas == nil {
			doc.Components.Schemas = openapi3.Schemas{}
		}
		if _, ok := doc.Components.Schemas[name]; !ok {
			if schemaRef, err := openapi3gen.NewSchemaRefForValue(*new(M), nil); err == nil {
				addSchemaTitleDesc[M](schemaRef)
				doc.Components.Schemas[name] = schemaRef
			}
		}
		docMutex.Unlock()
	}

	if !model.IsModelEmpty[REQ]() {
		typ := reflect.TypeOf(*new(M))
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		name := typ.Name()

		docMutex.Lock()
		if doc.Components.RequestBodies == nil {
			doc.Components.RequestBodies = openapi3.RequestBodies{}
		}
		if _, ok := doc.Components.RequestBodies[reqKey]; !ok && reqSchemaRef != nil {
			addSchemaTitleDesc[REQ](reqSchemaRef)
			setupExample(reqSchemaRef)
			setupBatchExample(reqSchemaRef)
			doc.Components.RequestBodies[reqKey] = &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: fmt.Sprintf("%s Payload", name),
					Required:    !model.IsModelEmpty[REQ](),
					Content:     openapi3.NewContentWithJSONSchemaRef(reqSchemaRef),
				},
			}

			// if schemaRef, err := openapi3gen.NewSchemaRefForValue(*new(REQ), nil); err == nil {
			// 	setupExample(schemaRef)
			// 	addSchemaTitleDesc[REQ](schemaRef)
			// 	doc.Components.RequestBodies[reqKey] = &openapi3.RequestBodyRef{
			// 		Value: &openapi3.RequestBody{
			// 			Description: fmt.Sprintf("%s payload", name),
			// 			Required:    !model.IsModelEmpty[REQ](),
			// 			Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef),
			// 		},
			// 	}
			// }
		}
		docMutex.Unlock()
	}

	if !model.IsModelEmpty[RSP]() {
		typ := reflect.TypeOf(*new(M))
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		name := typ.Name()

		docMutex.Lock()
		if doc.Components.Responses == nil {
			doc.Components.Responses = openapi3.ResponseBodies{}
		}
		if _, ok := doc.Components.Responses[rspKey]; !ok && rspSchemaRef != nil {
			doc.Components.Responses[rspKey] = &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: util.ValueOf(fmt.Sprintf("%s Response", name)),
					Content:     openapi3.NewContentWithJSONSchemaRef(rspSchemaRef),
				},
			}
			// if schemaRef, err := openapi3gen.NewSchemaRefForValue(*new(RSP), nil); err == nil {
			// 	addSchemaTitleDesc[RSP](schemaRef)
			// 	doc.Components.Responses[rspKey] = &openapi3.ResponseRef{
			// 		Value: &openapi3.Response{
			// 			Description: util.ValueOf(fmt.Sprintf("%s result", name)),
			// 			Content:     openapi3.NewContentWithJSONSchemaRef(schemaRef),
			// 		},
			// 	}
			// }
		}
		docMutex.Unlock()
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
//
// NOTE: 结构体字段必须有 json tag, 否则 schemaRef.Value.Properties 中不会带有这些字段
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

func setupBatchExample(schemaRef *openapi3.SchemaRef) {
	if schemaRef == nil || schemaRef.Value == nil {
		return
	}

	props := schemaRef.Value.Properties
	for k, v := range props {
		if k == "items" && v.Value != nil && v.Value.Type.Is(openapi3.TypeArray) {
			if v.Value.Items != nil && v.Value.Items.Value != nil {
				// 为数组中的单个元素创建 example
				example := make(map[string]any)
				for propName, propRef := range v.Value.Items.Value.Properties {
					if propName == "created_at" || propName == "created_by" || propName == "updated_at" || propName == "updated_by" || propName == "id" {
						continue
					}

					if propRef.Value == nil || propRef.Value.Type == nil {
						continue
					}

					switch {
					case propRef.Value.Type.Is(openapi3.TypeString):
						example[propName] = "string"
					case propRef.Value.Type.Is(openapi3.TypeInteger):
						example[propName] = 0
					case propRef.Value.Type.Is(openapi3.TypeNumber):
						example[propName] = 0.0
					case propRef.Value.Type.Is(openapi3.TypeBoolean):
						example[propName] = false
					case propRef.Value.Type.Is(openapi3.TypeArray):
						example[propName] = []any{}
					case propRef.Value.Type.Is(openapi3.TypeObject):
						example[propName] = map[string]any{}
					default:
						example[propName] = nil
					}
				}

				// 设置单个 item 的 example
				v.Value.Items.Value.Example = example

				// 设置整个 batch request 的 example
				schemaRef.Value.Example = map[string]any{
					"items": []map[string]any{example},
				}
			}
		}
	}
}

// removeFieldsFromRequestBody 从单个 CRUD 操作的 RequestBody 中移除指定字段
func removeFieldsFromRequestBody(op *openapi3.Operation, fieldsToRemove ...string) {
	if op == nil || op.RequestBody == nil {
		return
	}

	// 创建一个 map 方便查找
	removeMap := make(map[string]bool)
	for _, field := range fieldsToRemove {
		removeMap[field] = true
	}

	// 如果默认没有传入要移除的字段，使用默认值
	if len(fieldsToRemove) == 0 {
		removeMap = removeFieldMap
	}

	// 处理 RequestBodyRef
	var requestBody *openapi3.RequestBody

	if op.RequestBody.Ref != "" {
		// 如果是引用，需要从 components 中获取实际的 RequestBody
		// 这里需要访问全局的 doc 对象
		if doc.Components.RequestBodies != nil {
			refKey := strings.TrimPrefix(op.RequestBody.Ref, "#/components/requestBodies/")
			if rb, exists := doc.Components.RequestBodies[refKey]; exists && rb.Value != nil {
				requestBody = rb.Value
			}
		}
	} else if op.RequestBody.Value != nil {
		requestBody = op.RequestBody.Value
	}

	if requestBody == nil || requestBody.Content == nil {
		return
	}

	// 处理每个 content type
	for contentType, mediaType := range requestBody.Content {
		if mediaType.Schema == nil || mediaType.Schema.Value == nil {
			continue
		}

		schema := mediaType.Schema.Value

		// 移除 properties 中的字段
		if schema.Properties != nil {
			for field := range removeMap {
				delete(schema.Properties, field)
			}
		}

		// 移除 required 中的字段
		if len(schema.Required) > 0 {
			newRequired := []string{}
			for _, req := range schema.Required {
				if !removeMap[req] {
					newRequired = append(newRequired, req)
				}
			}
			schema.Required = newRequired
		}

		// 处理 example
		if schema.Example != nil {
			if exampleMap, ok := schema.Example.(map[string]any); ok {
				for field := range removeMap {
					delete(exampleMap, field)
				}
			}
		}

		// 更新 content
		requestBody.Content[contentType] = mediaType
	}
}

// removeFieldsFromBatchRequestBody 从批量 CRUD 操作的 RequestBody 中移除指定字段
func removeFieldsFromBatchRequestBody(op *openapi3.Operation, fieldsToRemove ...string) {
	if op == nil || op.RequestBody == nil {
		return
	}

	// 创建一个 map 方便查找
	removeMap := make(map[string]bool)
	for _, field := range fieldsToRemove {
		removeMap[field] = true
	}

	// 如果默认没有传入要移除的字段，使用默认值
	if len(fieldsToRemove) == 0 {
		removeMap = removeFieldMap
	}

	// 处理 RequestBodyRef
	var requestBody *openapi3.RequestBody

	if op.RequestBody.Ref != "" {
		// 如果是引用，需要从 components 中获取实际的 RequestBody
		if doc.Components.RequestBodies != nil {
			refKey := strings.TrimPrefix(op.RequestBody.Ref, "#/components/requestBodies/")
			if rb, exists := doc.Components.RequestBodies[refKey]; exists && rb.Value != nil {
				requestBody = rb.Value
			}
		}
	} else if op.RequestBody.Value != nil {
		requestBody = op.RequestBody.Value
	}

	if requestBody == nil || requestBody.Content == nil {
		return
	}

	// 处理每个 content type
	for contentType, mediaType := range requestBody.Content {
		if mediaType.Schema == nil || mediaType.Schema.Value == nil {
			continue
		}

		schema := mediaType.Schema.Value

		// 对于批量操作，需要处理 items 数组
		if schema.Properties != nil {
			if itemsProp, exists := schema.Properties["items"]; exists {
				if itemsProp.Value != nil && itemsProp.Value.Items != nil && itemsProp.Value.Items.Value != nil {
					itemSchema := itemsProp.Value.Items.Value

					// 移除 items 中每个元素的字段
					if itemSchema.Properties != nil {
						for field := range removeMap {
							delete(itemSchema.Properties, field)
						}
					}

					// 移除 required 中的字段
					if len(itemSchema.Required) > 0 {
						newRequired := []string{}
						for _, req := range itemSchema.Required {
							if !removeMap[req] {
								newRequired = append(newRequired, req)
							}
						}
						itemSchema.Required = newRequired
					}

					// 处理 items 的 example
					if itemSchema.Example != nil {
						if exampleMap, ok := itemSchema.Example.(map[string]any); ok {
							for field := range removeMap {
								delete(exampleMap, field)
							}
						}
					}
				}
			}
		}

		// 处理整个 batch request 的 example
		if schema.Example != nil {
			if exampleMap, ok := schema.Example.(map[string]any); ok {
				if items, exists := exampleMap["items"]; exists {
					if itemsArray, ok := items.([]map[string]any); ok {
						for _, item := range itemsArray {
							for field := range removeMap {
								delete(item, field)
							}
						}
					} else if itemsArray, ok := items.([]any); ok {
						for i, item := range itemsArray {
							if itemMap, ok := item.(map[string]any); ok {
								for field := range removeMap {
									delete(itemMap, field)
								}
								itemsArray[i] = itemMap
							}
						}
					}
				}
			}
		}

		// 更新 content
		requestBody.Content[contentType] = mediaType
	}
}

// 辅助函数：直接处理 schema，可以被上面两个函数调用
func removeFieldsFromSchema(schema *openapi3.Schema, fieldsToRemove map[string]bool) {
	if schema == nil {
		return
	}

	// 移除 properties
	if schema.Properties != nil {
		for field := range fieldsToRemove {
			delete(schema.Properties, field)
		}
	}

	// 移除 required
	if len(schema.Required) > 0 {
		newRequired := []string{}
		for _, req := range schema.Required {
			if !fieldsToRemove[req] {
				newRequired = append(newRequired, req)
			}
		}
		schema.Required = newRequired
	}

	// 处理 example
	if schema.Example != nil {
		if exampleMap, ok := schema.Example.(map[string]any); ok {
			for field := range fieldsToRemove {
				delete(exampleMap, field)
			}
		}
	}
}

// func setupBatchExample(schemaRef *openapi3.SchemaRef) {
// 	if schemaRef == nil {
// 		return
// 	}
// 	if schemaRef.Value == nil {
// 		schemaRef.Value = new(openapi3.Schema)
// 	}
// 	props := schemaRef.Value.Properties
// 	for k, v := range props {
// 		if k == "items" && v.Value.Type.Is(openapi3.TypeArray) {
// 			example := make(map[string]any)
// 			for k, v := range v.Value.Items.Value.Properties {
// 				if k == "created_at" || k == "created_by" || k == "updated_at" || k == "updated_by" {
// 					continue
// 				}
// 				if v.Value == nil || v.Value.Type == nil {
// 					continue
// 				}
// 				if v.Value.Type.Is(openapi3.TypeString) {
// 					example[k] = "string"
// 				}
// 				if v.Value.Type.Is(openapi3.TypeInteger) {
// 					example[k] = 0
// 				}
// 				if v.Value.Type.Is(openapi3.TypeNumber) {
// 					example[k] = 0.0
// 				}
// 				if v.Value.Type.Is(openapi3.TypeBoolean) {
// 					example[k] = false
// 				}
// 				if v.Value.Type.Is(openapi3.TypeArray) {
// 					example[k] = []any{}
// 				}
// 				if v.Value.Type.Is(openapi3.TypeObject) {
// 					example[k] = map[string]any{}
// 				}
// 				if v.Value.Type.Is(openapi3.TypeNull) {
// 					example[k] = nil
// 				}
// 			}
// 			v.Value.Items.Value.Example = example
// 		}
// 	}
// }

func addHeaderParameters(op *openapi3.Operation) {
	headers := []*openapi3.ParameterRef{
		// // Mybe used in the future, DO NOT DELETE it.
		// {
		// 	Value: &openapi3.Parameter{
		// 		In:          "header",
		// 		Name:        "Authorization",
		// 		Description: "Authentication token (e.g. Bearer <token>)",
		// 		Required:    false,
		// 		Schema: &openapi3.SchemaRef{
		// 			Value: &openapi3.Schema{
		// 				Type: &openapi3.Types{openapi3.TypeString},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	Value: &openapi3.Parameter{
		// 		In:          "header",
		// 		Name:        "X-Request-ID",
		// 		Description: "Optional request ID for tracing",
		// 		Required:    false,
		// 		Schema: &openapi3.SchemaRef{
		// 			Value: &openapi3.Schema{
		// 				Type: &openapi3.Types{openapi3.TypeString},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	Value: &openapi3.Parameter{
		// 		In:          "header",
		// 		Name:        "X-Client-Version",
		// 		Description: "Client version (e.g. v1.2.3)",
		// 		Required:    false,
		// 		Schema: &openapi3.SchemaRef{
		// 			Value: &openapi3.Schema{
		// 				Type: &openapi3.Types{openapi3.TypeString},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	Value: &openapi3.Parameter{
		// 		In:          "header",
		// 		Name:        "Accept-Language",
		// 		Description: "Preferred language (e.g. zh-CN, en-US)",
		// 		Required:    false,
		// 		Schema: &openapi3.SchemaRef{
		// 			Value: &openapi3.Schema{
		// 				Type: &openapi3.Types{openapi3.TypeString},
		// 			},
		// 		},
		// 	},
		// },
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
func addSchemaTitleDesc[T any](schemaRef *openapi3.SchemaRef) {
	if schemaRef == nil || schemaRef.Value == nil || schemaRef.Value.Properties == nil {
		return
	}

	// Get model field descriptions
	modelInstance := *new(T)
	modelDocs := parseModelDocs(modelInstance)

	// Get field descriptions of model.Base (using cache)
	baseDocs := getBaseModelDocs()

	// Create a mapping from JSON property names to field descriptions
	propertyDescriptions := make(map[string]string)

	// Process model fields
	typ := reflect.TypeOf(*new(T)).Elem()
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

// addQueryParameters adds query parameters for List operation.
func addQueryParameters[M types.Model, REQ types.Request, RSP types.Response](op *openapi3.Operation) {
	// 只有使用默认的逻辑才支持通过结构体字段过滤
	if !model.AreTypesEqual[M, REQ, RSP]() {
		return
	}

	queries := make([]*openapi3.ParameterRef, 0)

	// Get model field descriptions
	modelInstance := *new(M)
	modelDocs := parseModelDocs(modelInstance)

	typ := reflect.TypeOf(*new(M)).Elem()
	for i := range typ.NumField() {
		field := typ.Field(i)
		// 只有增加了 schema 标签的字段才支持通过 request query 参数进行后端查询
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
	items := strings.Split(tag, `/`)
	if len(items) > 0 {
		tag = items[0]
	} else {
		tag = typ.Elem().Name()
	}
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

func newRequestBody(reqKey string) *openapi3.RequestBodyRef {
	return &openapi3.RequestBodyRef{
		Ref: "#/components/requestBodies/" + reqKey,
	}
}

func newResponses(status int, rspKey string) *openapi3.Responses {
	return openapi3.NewResponses(openapi3.WithStatus(status, &openapi3.ResponseRef{Ref: "#/components/responses/" + rspKey}))
}

// func NewResponses() *openapi3.Responses {
// 	if len(opts) == 0 {
// 		return NewResponses(WithName("default", NewResponse().WithDescription("")))
// 	}
// 	return NewResponses(openapi3.WithName())
// }

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

func registerCommonResponses() {
	if doc.Components.Responses == nil {
		doc.Components.Responses = openapi3.ResponseBodies{}
	}

	// 400 Bad Request
	doc.Components.Responses["BadRequest"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Bad Request - The request was invalid or cannot be served"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 400,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Invalid request parameters",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
					Example: map[string]any{
						"code":       400,
						"msg":        "Invalid request parameters",
						"request_id": "req_123456789",
						"data":       nil,
					},
				},
			}),
		},
	}

	// 401 Unauthorized
	doc.Components.Responses["Unauthorized"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Unauthorized - Authentication credentials were missing or incorrect"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 401,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Authentication required",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
					Example: map[string]any{
						"code":       401,
						"msg":        "Authentication required",
						"request_id": "req_123456789",
						"data":       nil,
					},
				},
			}),
		},
	}

	// 403 Forbidden
	doc.Components.Responses["Forbidden"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Forbidden - The request is understood, but access is not allowed"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 403,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Access denied",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
					Example: map[string]any{
						"code":       403,
						"msg":        "Access denied",
						"request_id": "req_123456789",
						"data":       nil,
					},
				},
			}),
		},
	}

	// 404 Not Found
	doc.Components.Responses["NotFound"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Not Found - The requested resource could not be found"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 404,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Resource not found",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
					Example: map[string]any{
						"code":       404,
						"msg":        "Resource not found",
						"request_id": "req_123456789",
						"data":       nil,
					},
				},
			}),
		},
	}

	// 409 Conflict
	doc.Components.Responses["Conflict"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Conflict - The request could not be completed due to a conflict with the current state"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 409,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Resource already exists",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
					Example: map[string]any{
						"code":       409,
						"msg":        "Resource already exists",
						"request_id": "req_123456789",
						"data":       nil,
					},
				},
			}),
		},
	}

	// 422 Unprocessable Entity
	doc.Components.Responses["UnprocessableEntity"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Unprocessable Entity - The request was well-formed but contains semantic errors"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 422,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Validation failed",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:        &openapi3.Types{openapi3.TypeObject},
								Nullable:    true,
								Description: "Validation errors detail",
								Example: map[string]any{
									"errors": []map[string]any{
										{
											"field":   "name",
											"message": "Name is required",
										},
									},
								},
							},
						},
					},
					Example: map[string]any{
						"code":       422,
						"msg":        "Validation failed",
						"request_id": "req_123456789",
						"data": map[string]any{
							"errors": []map[string]any{
								{
									"field":   "name",
									"message": "Name is required",
								},
							},
						},
					},
				},
			}),
		},
	}

	// 429 Too Many Requests
	doc.Components.Responses["TooManyRequests"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Too Many Requests - Rate limit exceeded"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 429,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Rate limit exceeded. Please try again later",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
					Example: map[string]any{
						"code":       429,
						"msg":        "Rate limit exceeded. Please try again later",
						"request_id": "req_123456789",
						"data":       nil,
					},
				},
			}),
		},
	}

	// 500 Internal Server Error
	doc.Components.Responses["InternalServerError"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Internal Server Error - An unexpected error occurred"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 500,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Internal server error",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
					Example: map[string]any{
						"code":       500,
						"msg":        "Internal server error",
						"request_id": "req_123456789",
						"data":       nil,
					},
				},
			}),
		},
	}

	// 502 Bad Gateway
	doc.Components.Responses["BadGateway"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Bad Gateway - Invalid response from upstream server"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 502,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Bad gateway",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
				},
			}),
		},
	}

	// 503 Service Unavailable
	doc.Components.Responses["ServiceUnavailable"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: util.ValueOf("Service Unavailable - The server is currently unable to handle the request"),
			Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeInteger},
								Example: 503,
							},
						},
						"msg": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "Service temporarily unavailable",
							},
						},
						"request_id": {
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{openapi3.TypeString},
								Example: "req_123456789",
							},
						},
						"data": {
							Value: &openapi3.Schema{
								Type:     &openapi3.Types{openapi3.TypeNull},
								Nullable: true,
							},
						},
					},
				},
			}),
		},
	}
}

func registerSecuritySchemes() {
	if doc.Components.SecuritySchemes == nil {
		doc.Components.SecuritySchemes = openapi3.SecuritySchemes{}
	}

	// Bearer Token
	doc.Components.SecuritySchemes["bearerAuth"] = &openapi3.SecuritySchemeRef{
		Value: &openapi3.SecurityScheme{
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "Enter JWT Bearer token",
		},
	}

	// API Key
	doc.Components.SecuritySchemes["apiKey"] = &openapi3.SecuritySchemeRef{
		Value: &openapi3.SecurityScheme{
			Type:        "apiKey",
			In:          "header",
			Name:        "X-API-Key",
			Description: "API Key authentication",
		},
	}
}

// 在文档级别应用安全要求
func applyGlobalSecurity() {
	doc.Security = openapi3.SecurityRequirements{
		{
			"bearerAuth": []string{},
		},
	}
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
