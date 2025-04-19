package openapigen

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

// 工具函数：字符串转指针
func strPtr(s string) *string {
	return &s
}

type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
	CreatedBy *string `json:"created_by,omitempty"`
	UpdatedBy *string `json:"updated_by,omitempty"`
	Remark    *string `json:"remark,omitempty"`
}

type UserInput struct {
	Name string `json:"name"`
}

type UserPatch struct {
	Name *string `json:"name,omitempty"`
}

func GenerateOpenAPI() (*openapi3.T, error) {
	doc := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "User API",
			Description: "API for managing users",
			Version:     "1.0.0",
		},
		Paths: openapi3.NewPaths(),
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{},
		},
	}

	// === 使用 openapi3gen 自动生成 Schema ===
	gen := openapi3gen.NewGenerator()

	userSchemaRef, err := gen.NewSchemaRefForValue(User{}, nil)
	if err != nil {
		return nil, err
	}
	userInputRef, err := gen.NewSchemaRefForValue(UserInput{}, nil)
	if err != nil {
		return nil, err
	}
	userPatchRef, err := gen.NewSchemaRefForValue(UserPatch{}, nil)
	if err != nil {
		return nil, err
	}

	// 注册组件（可选）
	doc.Components.Schemas["User"] = userSchemaRef
	doc.Components.Schemas["UserInput"] = userInputRef
	doc.Components.Schemas["UserPatch"] = userPatchRef

	// === 定义响应 ===
	getResponses := openapi3.NewResponses()
	getResponses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("List of users"),
			Content: openapi3.NewContentWithJSONSchema(
				openapi3.NewArraySchema().WithItems(userSchemaRef.Value),
			),
		},
	})

	postResponses := openapi3.NewResponses()
	postResponses.Set("201", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User created"),
			Content:     openapi3.NewContentWithJSONSchemaRef(userSchemaRef),
		},
	})

	doc.Paths.Set("/api/user", &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUsers",
			Summary:     "Get list of users",
			Responses:   getResponses,
		},
		Post: &openapi3.Operation{
			OperationID: "createUser",
			Summary:     "Create a new user",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "User to create",
					Required:    true,
					Content:     openapi3.NewContentWithJSONSchemaRef(userInputRef),
				},
			},
			Responses: postResponses,
		},
	})

	// 通用 path 参数 {id}
	userIdParam := &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			In:       "path",
			Name:     "id",
			Required: true,
			Schema:   openapi3.NewIntegerSchema().NewRef(),
		},
	}

	// GET /api/user/{id}
	getByIdResponses := openapi3.NewResponses()
	getByIdResponses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User found"),
			Content:     openapi3.NewContentWithJSONSchemaRef(userSchemaRef),
		},
	})
	getByIdResponses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User not found"),
		},
	})

	// PUT /api/user/{id}
	putResponses := openapi3.NewResponses()
	putResponses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User updated"),
			Content:     openapi3.NewContentWithJSONSchemaRef(userSchemaRef),
		},
	})
	putResponses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User not found"),
		},
	})

	// DELETE
	deleteResponses := openapi3.NewResponses()
	deleteResponses.Set("204", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User deleted"),
		},
	})
	deleteResponses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User not found"),
		},
	})

	// PATCH
	patchResponses := openapi3.NewResponses()
	patchResponses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User partially updated"),
			Content:     openapi3.NewContentWithJSONSchemaRef(userSchemaRef),
		},
	})
	patchResponses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: strPtr("User not found"),
		},
	})

	// 添加 /api/user/{id}
	doc.Paths.Set("/api/user/{id}", &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUserById",
			Summary:     "Retrieve a user by ID",
			Parameters:  []*openapi3.ParameterRef{userIdParam},
			Responses:   getByIdResponses,
		},
		Put: &openapi3.Operation{
			OperationID: "updateUser",
			Summary:     "Update a user by ID",
			Parameters:  []*openapi3.ParameterRef{userIdParam},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Updated user info",
					Required:    true,
					Content:     openapi3.NewContentWithJSONSchemaRef(userInputRef),
				},
			},
			Responses: putResponses,
		},
		Delete: &openapi3.Operation{
			OperationID: "deleteUser",
			Summary:     "Delete a user by ID",
			Parameters:  []*openapi3.ParameterRef{userIdParam},
			Responses:   deleteResponses,
		},
		Patch: &openapi3.Operation{
			OperationID: "patchUser",
			Summary:     "Partially update a user by ID",
			Parameters:  []*openapi3.ParameterRef{userIdParam},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Partial user update",
					Required:    true,
					Content:     openapi3.NewContentWithJSONSchemaRef(userPatchRef),
				},
			},
			Responses: patchResponses,
		},
	})

	// 写入 api.json 文件
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}

	return doc, os.WriteFile("api.json", data, 0o644)
}

// DocumentHandler returns an http.Handler that serves the OpenAPI document
func DocumentHandler(doc *openapi3.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.MarshalIndent(doc, "", "  ")
		w.Write(data)
	})
}
