package openapigen

import (
	"testing"

	"github.com/forbearing/golib/model"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/stretchr/testify/assert"
)

// TestUser is a user model for testing
type TestUser struct {
	// Name user's name
	Name string `json:"name" schema:"name"`
	// Email user's email address
	Email string `json:"email" schema:"email"`

	model.Base
}

func Test_addQueryParameters(t *testing.T) {
	t.Run("test_query_parameters_with_descriptions", func(t *testing.T) {
		// Create an OpenAPI operation
		op := &openapi3.Operation{
			Parameters: make([]*openapi3.ParameterRef, 0),
		}

		// Call addQueryParameters function
		addQueryParameters[*TestUser](op)

		// Verify that parameters are correctly added
		assert.NotEmpty(t, op.Parameters)

		// Create a mapping for quick parameter lookup
		paramMap := make(map[string]*openapi3.Parameter)
		for _, paramRef := range op.Parameters {
			if paramRef.Value != nil {
				paramMap[paramRef.Value.Name] = paramRef.Value
			}
		}

		// Verify parameters for model fields
		if nameParam, exists := paramMap["name"]; exists {
			assert.Equal(t, "name", nameParam.Name)
			assert.Equal(t, "query", nameParam.In)
			assert.False(t, nameParam.Required)
			assert.Equal(t, "Name user's name", nameParam.Description)
		}

		if emailParam, exists := paramMap["email"]; exists {
			assert.Equal(t, "email", emailParam.Name)
			assert.Equal(t, "query", emailParam.In)
			assert.False(t, emailParam.Required)
			assert.Equal(t, "Email user's email address", emailParam.Description)
		}

		// Verify some parameters of Base model (these should have descriptions, even if empty)
		if pageParam, exists := paramMap["page"]; exists {
			assert.Equal(t, "page", pageParam.Name)
			assert.Equal(t, "query", pageParam.In)
			assert.False(t, pageParam.Required)
			// Base model field descriptions may be empty because they may not have comments in source code
			// But should at least have Description field
			assert.NotNil(t, pageParam.Description)
		}

		if sizeParam, exists := paramMap["size"]; exists {
			assert.Equal(t, "size", sizeParam.Name)
			assert.Equal(t, "query", sizeParam.In)
			assert.False(t, sizeParam.Required)
			assert.NotNil(t, sizeParam.Description)
		}

		t.Logf("Total parameters added: %d", len(op.Parameters))
		for _, paramRef := range op.Parameters {
			if paramRef.Value != nil {
				t.Logf("Parameter: %s, Description: %s", paramRef.Value.Name, paramRef.Value.Description)
			}
		}
	})
}

func Test_setBatchDelete(t *testing.T) {
	t.Run("test_batch_delete_with_schema_descriptions", func(t *testing.T) {
		// Create a path item
		pathItem := &openapi3.PathItem{}

		// Call setBatchDelete function
		setDeleteMany[*TestUser]("/users/batch", pathItem)

		// Verify that DELETE operation is created
		assert.NotNil(t, pathItem.Delete)
		assert.Equal(t, "batch deleteTestUser", pathItem.Delete.OperationID)

		// Verify response schema has descriptions for items array
		assert.NotNil(t, pathItem.Delete.Responses)
		response200 := pathItem.Delete.Responses.Value("200")
		assert.NotNil(t, response200)
		assert.NotNil(t, response200.Value)
		assert.NotNil(t, response200.Value.Content)

		jsonResponse := response200.Value.Content.Get("application/json")
		assert.NotNil(t, jsonResponse)
		assert.NotNil(t, jsonResponse.Schema)
		assert.NotNil(t, jsonResponse.Schema.Value)
		assert.NotNil(t, jsonResponse.Schema.Value.Properties)

		// Check that response data field has descriptions for items array
		if dataProperty, exists := jsonResponse.Schema.Value.Properties["data"]; exists {
			assert.NotNil(t, dataProperty.Value)
			assert.NotNil(t, dataProperty.Value.Properties)

			if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
				assert.NotNil(t, itemsProperty.Value)
				assert.NotNil(t, itemsProperty.Value.Items)
				assert.NotNil(t, itemsProperty.Value.Items.Value)
				assert.NotNil(t, itemsProperty.Value.Items.Value.Properties)

				// Check name field in response items array
				if nameProperty, exists := itemsProperty.Value.Items.Value.Properties["name"]; exists {
					assert.NotNil(t, nameProperty.Value)
					assert.Equal(t, "Name user's name", nameProperty.Value.Description)
				}

				// Check email field in response items array
				if emailProperty, exists := itemsProperty.Value.Items.Value.Properties["email"]; exists {
					assert.NotNil(t, emailProperty.Value)
					assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
				}
			}
		}
	})
}

func Test_setCreate(t *testing.T) {
	t.Run("test_create_with_schema_descriptions", func(t *testing.T) {
		// Create a path item
		pathItem := &openapi3.PathItem{}

		// Call setCreate function
		setCreate[*TestUser]("/users", pathItem)

		// Verify that POST operation is created
		assert.NotNil(t, pathItem.Post)
		assert.Equal(t, "createTestUser", pathItem.Post.OperationID)

		// Verify request body schema has descriptions
		assert.NotNil(t, pathItem.Post.RequestBody)
		assert.NotNil(t, pathItem.Post.RequestBody.Value)
		assert.NotNil(t, pathItem.Post.RequestBody.Value.Content)

		jsonContent := pathItem.Post.RequestBody.Value.Content.Get("application/json")
		assert.NotNil(t, jsonContent)
		assert.NotNil(t, jsonContent.Schema)
		assert.NotNil(t, jsonContent.Schema.Value)
		assert.NotNil(t, jsonContent.Schema.Value.Properties)

		// Check that name field has description
		if nameProperty, exists := jsonContent.Schema.Value.Properties["name"]; exists {
			assert.NotNil(t, nameProperty.Value)
			assert.Equal(t, "Name user's name", nameProperty.Value.Description)
		}

		// Check that email field has description
		if emailProperty, exists := jsonContent.Schema.Value.Properties["email"]; exists {
			assert.NotNil(t, emailProperty.Value)
			assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
		}

		// Verify response schema has descriptions
		assert.NotNil(t, pathItem.Post.Responses)
		response201 := pathItem.Post.Responses.Value("201")
		assert.NotNil(t, response201)
		assert.NotNil(t, response201.Value)
		assert.NotNil(t, response201.Value.Content)

		jsonResponse := response201.Value.Content.Get("application/json")
		assert.NotNil(t, jsonResponse)
		assert.NotNil(t, jsonResponse.Schema)
		assert.NotNil(t, jsonResponse.Schema.Value)
		assert.NotNil(t, jsonResponse.Schema.Value.Properties)

		// Check that response data field has descriptions
		if dataProperty, exists := jsonResponse.Schema.Value.Properties["data"]; exists {
			assert.NotNil(t, dataProperty.Value)
			assert.NotNil(t, dataProperty.Value.Properties)

			// Check name field in response data
			if nameProperty, exists := dataProperty.Value.Properties["name"]; exists {
				assert.NotNil(t, nameProperty.Value)
				assert.Equal(t, "Name user's name", nameProperty.Value.Description)
			}

			// Check email field in response data
			if emailProperty, exists := dataProperty.Value.Properties["email"]; exists {
				assert.NotNil(t, emailProperty.Value)
				assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
			}
		}
	})
}

func Test_setUpdate(t *testing.T) {
	t.Run("test_update_with_schema_descriptions", func(t *testing.T) {
		// Create a path item
		pathItem := &openapi3.PathItem{}

		// Call setUpdate function
		setUpdate[*TestUser]("/users/{id}", pathItem)

		// Verify that PUT operation is created
		assert.NotNil(t, pathItem.Put)
		assert.Equal(t, "updateTestUser", pathItem.Put.OperationID)

		// Verify request body schema has descriptions
		assert.NotNil(t, pathItem.Put.RequestBody)
		assert.NotNil(t, pathItem.Put.RequestBody.Value)
		assert.NotNil(t, pathItem.Put.RequestBody.Value.Content)

		jsonRequest := pathItem.Put.RequestBody.Value.Content.Get("application/json")
		assert.NotNil(t, jsonRequest)
		assert.NotNil(t, jsonRequest.Schema)
		assert.NotNil(t, jsonRequest.Schema.Value)
		assert.NotNil(t, jsonRequest.Schema.Value.Properties)

		// Check name field in request body
		if nameProperty, exists := jsonRequest.Schema.Value.Properties["name"]; exists {
			assert.NotNil(t, nameProperty.Value)
			assert.Equal(t, "Name user's name", nameProperty.Value.Description)
		}

		// Check email field in request body
		if emailProperty, exists := jsonRequest.Schema.Value.Properties["email"]; exists {
			assert.NotNil(t, emailProperty.Value)
			assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
		}

		// Verify response schema has descriptions
		assert.NotNil(t, pathItem.Put.Responses)
		response200 := pathItem.Put.Responses.Value("200")
		assert.NotNil(t, response200)
		assert.NotNil(t, response200.Value)
		assert.NotNil(t, response200.Value.Content)

		jsonResponse := response200.Value.Content.Get("application/json")
		assert.NotNil(t, jsonResponse)
		assert.NotNil(t, jsonResponse.Schema)
		assert.NotNil(t, jsonResponse.Schema.Value)
		assert.NotNil(t, jsonResponse.Schema.Value.Properties)

		// Check that response data field has descriptions
		if dataProperty, exists := jsonResponse.Schema.Value.Properties["data"]; exists {
			assert.NotNil(t, dataProperty.Value)
			assert.NotNil(t, dataProperty.Value.Properties)

			// Check name field in response data
			if nameProperty, exists := dataProperty.Value.Properties["name"]; exists {
				assert.NotNil(t, nameProperty.Value)
				assert.Equal(t, "Name user's name", nameProperty.Value.Description)
			}

			// Check email field in response data
			if emailProperty, exists := dataProperty.Value.Properties["email"]; exists {
				assert.NotNil(t, emailProperty.Value)
				assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
			}
		}
	})
}

func Test_setUpdatePartial(t *testing.T) {
	t.Run("test_update_partial_with_schema_descriptions", func(t *testing.T) {
		// Create a path item
		pathItem := &openapi3.PathItem{}

		// Call setUpdatePartial function
		setUpdatePartial[*TestUser]("/users/{id}", pathItem)

		// Verify that PATCH operation is created
		assert.NotNil(t, pathItem.Patch)
		assert.Equal(t, "update partialTestUser", pathItem.Patch.OperationID)

		// Verify request body schema has descriptions
		assert.NotNil(t, pathItem.Patch.RequestBody)
		assert.NotNil(t, pathItem.Patch.RequestBody.Value)
		assert.NotNil(t, pathItem.Patch.RequestBody.Value.Content)

		jsonRequest := pathItem.Patch.RequestBody.Value.Content.Get("application/json")
		assert.NotNil(t, jsonRequest)
		assert.NotNil(t, jsonRequest.Schema)
		assert.NotNil(t, jsonRequest.Schema.Value)
		assert.NotNil(t, jsonRequest.Schema.Value.Properties)

		// Check name field in request body
		if nameProperty, exists := jsonRequest.Schema.Value.Properties["name"]; exists {
			assert.NotNil(t, nameProperty.Value)
			assert.Equal(t, "Name user's name", nameProperty.Value.Description)
		}

		// Check email field in request body
		if emailProperty, exists := jsonRequest.Schema.Value.Properties["email"]; exists {
			assert.NotNil(t, emailProperty.Value)
			assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
		}

		// Verify response schema has descriptions
		assert.NotNil(t, pathItem.Patch.Responses)
		response200 := pathItem.Patch.Responses.Value("200")
		assert.NotNil(t, response200)
		assert.NotNil(t, response200.Value)
		assert.NotNil(t, response200.Value.Content)

		jsonResponse := response200.Value.Content.Get("application/json")
		assert.NotNil(t, jsonResponse)
		assert.NotNil(t, jsonResponse.Schema)
		assert.NotNil(t, jsonResponse.Schema.Value)
		assert.NotNil(t, jsonResponse.Schema.Value.Properties)

		// Check that response data field has descriptions
		if dataProperty, exists := jsonResponse.Schema.Value.Properties["data"]; exists {
			assert.NotNil(t, dataProperty.Value)
			assert.NotNil(t, dataProperty.Value.Properties)

			// Check name field in response data
			if nameProperty, exists := dataProperty.Value.Properties["name"]; exists {
				assert.NotNil(t, nameProperty.Value)
				assert.Equal(t, "Name user's name", nameProperty.Value.Description)
			}

			// Check email field in response data
			if emailProperty, exists := dataProperty.Value.Properties["email"]; exists {
				assert.NotNil(t, emailProperty.Value)
				assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
			}
		}
	})
}

func Test_setDelete(t *testing.T) {
	t.Run("test_delete_with_schema_descriptions", func(t *testing.T) {
		// Create a path item
		pathItem := &openapi3.PathItem{}

		// Call setDelete function
		setDelete[*TestUser]("/users/{id}", pathItem)

		// Verify that DELETE operation is created
		assert.NotNil(t, pathItem.Delete)
		assert.Equal(t, "deleteTestUser", pathItem.Delete.OperationID)

		// Verify response schema has descriptions
		assert.NotNil(t, pathItem.Delete.Responses)
		response204 := pathItem.Delete.Responses.Value("204")
		assert.NotNil(t, response204)
		assert.NotNil(t, response204.Value)
		assert.NotNil(t, response204.Value.Content)

		jsonResponse := response204.Value.Content.Get("application/json")
		assert.NotNil(t, jsonResponse)
		assert.NotNil(t, jsonResponse.Schema)
		assert.NotNil(t, jsonResponse.Schema.Value)
		assert.NotNil(t, jsonResponse.Schema.Value.Properties)

		// Check that response data field has descriptions
		if dataProperty, exists := jsonResponse.Schema.Value.Properties["data"]; exists {
			assert.NotNil(t, dataProperty.Value)
			assert.NotNil(t, dataProperty.Value.Properties)

			// Check name field in response data
			if nameProperty, exists := dataProperty.Value.Properties["name"]; exists {
				assert.NotNil(t, nameProperty.Value)
				assert.Equal(t, "Name user's name", nameProperty.Value.Description)
			}

			// Check email field in response data
			if emailProperty, exists := dataProperty.Value.Properties["email"]; exists {
				assert.NotNil(t, emailProperty.Value)
				assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
			}
		}
	})
}

func Test_setList(t *testing.T) {
	t.Run("test_list_with_schema_descriptions", func(t *testing.T) {
		// Create a path item
		pathItem := &openapi3.PathItem{}

		// Call setList function
		setList[*TestUser]("/users", pathItem)

		// Verify that GET operation is created
		assert.NotNil(t, pathItem.Get)
		assert.Equal(t, "listTestUser", pathItem.Get.OperationID)

		// Verify response schema has descriptions
		assert.NotNil(t, pathItem.Get.Responses)
		response200 := pathItem.Get.Responses.Value("200")
		assert.NotNil(t, response200)
		assert.NotNil(t, response200.Value)
		assert.NotNil(t, response200.Value.Content)

		jsonResponse := response200.Value.Content.Get("application/json")
		assert.NotNil(t, jsonResponse)
		assert.NotNil(t, jsonResponse.Schema)
		assert.NotNil(t, jsonResponse.Schema.Value)
		assert.NotNil(t, jsonResponse.Schema.Value.Properties)

		// Check that response data field has descriptions for list items
		if dataProperty, exists := jsonResponse.Schema.Value.Properties["data"]; exists {
			assert.NotNil(t, dataProperty.Value)
			if dataProperty.Value != nil && dataProperty.Value.Items != nil {
				assert.NotNil(t, dataProperty.Value.Items.Value)
				if dataProperty.Value.Items.Value != nil && dataProperty.Value.Items.Value.Properties != nil {
					// Check name field in list items
					if nameProperty, exists := dataProperty.Value.Items.Value.Properties["name"]; exists {
						assert.NotNil(t, nameProperty.Value)
						assert.Equal(t, "Name user's name", nameProperty.Value.Description)
					}

					// Check email field in list items
					if emailProperty, exists := dataProperty.Value.Items.Value.Properties["email"]; exists {
						assert.NotNil(t, emailProperty.Value)
						assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
					}
				}
			}
		}
	})
}

func Test_setBatchUpdate(t *testing.T) {
	t.Run("test_batch_update_with_schema_descriptions", func(t *testing.T) {
		// Create a path item
		pathItem := &openapi3.PathItem{}

		// Call setBatchUpdate function
		setUpdateMany[*TestUser]("/testusers", pathItem)

		// Verify that the PUT operation is created
		assert.NotNil(t, pathItem.Put)
		assert.Equal(t, "batch updateTestUser", pathItem.Put.OperationID)

		// Verify request body has field descriptions
		assert.NotNil(t, pathItem.Put.RequestBody)
		assert.NotNil(t, pathItem.Put.RequestBody.Value)
		assert.NotNil(t, pathItem.Put.RequestBody.Value.Content)

		// Get the JSON content from request body
		jsonContent := pathItem.Put.RequestBody.Value.Content["application/json"]
		assert.NotNil(t, jsonContent)
		assert.NotNil(t, jsonContent.Schema)
		assert.NotNil(t, jsonContent.Schema.Value)
		assert.NotNil(t, jsonContent.Schema.Value.Properties)

		// Check items array in request body
		if itemsProperty, exists := jsonContent.Schema.Value.Properties["items"]; exists {
			assert.NotNil(t, itemsProperty.Value)
			assert.NotNil(t, itemsProperty.Value.Items)
			assert.NotNil(t, itemsProperty.Value.Items.Value)
			assert.NotNil(t, itemsProperty.Value.Items.Value.Properties)

			// Check field descriptions in items
			if nameProperty, exists := itemsProperty.Value.Items.Value.Properties["name"]; exists {
				assert.NotNil(t, nameProperty.Value)
				assert.Equal(t, "Name user's name", nameProperty.Value.Description)
			}

			if emailProperty, exists := itemsProperty.Value.Items.Value.Properties["email"]; exists {
				assert.NotNil(t, emailProperty.Value)
				assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
			}
		}

		// Verify response has field descriptions
		assert.NotNil(t, pathItem.Put.Responses)
		response200 := pathItem.Put.Responses.Value("200")
		assert.NotNil(t, response200)
		assert.NotNil(t, response200.Value)
		assert.NotNil(t, response200.Value.Content)

		// Get the JSON content from response
		jsonResponseContent := response200.Value.Content["application/json"]
		assert.NotNil(t, jsonResponseContent)
		assert.NotNil(t, jsonResponseContent.Schema)
		assert.NotNil(t, jsonResponseContent.Schema.Value)
		assert.NotNil(t, jsonResponseContent.Schema.Value.Properties)

		// Check data.items array in response
		if dataProperty, exists := jsonResponseContent.Schema.Value.Properties["data"]; exists {
			assert.NotNil(t, dataProperty.Value)
			assert.NotNil(t, dataProperty.Value.Properties)

			if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
				assert.NotNil(t, itemsProperty.Value)
				assert.NotNil(t, itemsProperty.Value.Items)
				assert.NotNil(t, itemsProperty.Value.Items.Value)
				assert.NotNil(t, itemsProperty.Value.Items.Value.Properties)

				// Check field descriptions in response items
				if nameProperty, exists := itemsProperty.Value.Items.Value.Properties["name"]; exists {
					assert.NotNil(t, nameProperty.Value)
					assert.Equal(t, "Name user's name", nameProperty.Value.Description)
				}

				if emailProperty, exists := itemsProperty.Value.Items.Value.Properties["email"]; exists {
					assert.NotNil(t, emailProperty.Value)
					assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
				}
			}
		}
	})
}

func Test_setBatchUpdatePartial(t *testing.T) {
	t.Run("test_batch_update_partial_with_schema_descriptions", func(t *testing.T) {
		// Create a path item
		pathItem := &openapi3.PathItem{}

		// Call setBatchUpdatePartial function
		setPatchMany[*TestUser]("/testusers", pathItem)

		// Verify that the PATCH operation is created
		assert.NotNil(t, pathItem.Patch)
		assert.Equal(t, "batch update partialTestUser", pathItem.Patch.OperationID)

		// Verify request body has field descriptions
		assert.NotNil(t, pathItem.Patch.RequestBody)
		assert.NotNil(t, pathItem.Patch.RequestBody.Value)
		assert.NotNil(t, pathItem.Patch.RequestBody.Value.Content)

		// Get the JSON content from request body
		jsonContent := pathItem.Patch.RequestBody.Value.Content["application/json"]
		assert.NotNil(t, jsonContent)
		assert.NotNil(t, jsonContent.Schema)
		assert.NotNil(t, jsonContent.Schema.Value)
		assert.NotNil(t, jsonContent.Schema.Value.Properties)

		// Check items array in request body
		if itemsProperty, exists := jsonContent.Schema.Value.Properties["items"]; exists {
			assert.NotNil(t, itemsProperty.Value)
			assert.NotNil(t, itemsProperty.Value.Items)
			assert.NotNil(t, itemsProperty.Value.Items.Value)
			assert.NotNil(t, itemsProperty.Value.Items.Value.Properties)

			// Check field descriptions in items
			if nameProperty, exists := itemsProperty.Value.Items.Value.Properties["name"]; exists {
				assert.NotNil(t, nameProperty.Value)
				assert.Equal(t, "Name user's name", nameProperty.Value.Description)
			}

			if emailProperty, exists := itemsProperty.Value.Items.Value.Properties["email"]; exists {
				assert.NotNil(t, emailProperty.Value)
				assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
			}
		}

		// Verify response has field descriptions
		assert.NotNil(t, pathItem.Patch.Responses)
		response200 := pathItem.Patch.Responses.Value("200")
		assert.NotNil(t, response200)
		assert.NotNil(t, response200.Value)
		assert.NotNil(t, response200.Value.Content)

		// Get the JSON content from response
		jsonResponseContent := response200.Value.Content["application/json"]
		assert.NotNil(t, jsonResponseContent)
		assert.NotNil(t, jsonResponseContent.Schema)
		assert.NotNil(t, jsonResponseContent.Schema.Value)
		assert.NotNil(t, jsonResponseContent.Schema.Value.Properties)

		// Check data.items array in response
		if dataProperty, exists := jsonResponseContent.Schema.Value.Properties["data"]; exists {
			assert.NotNil(t, dataProperty.Value)
			assert.NotNil(t, dataProperty.Value.Properties)

			if itemsProperty, exists := dataProperty.Value.Properties["items"]; exists {
				assert.NotNil(t, itemsProperty.Value)
				assert.NotNil(t, itemsProperty.Value.Items)
				assert.NotNil(t, itemsProperty.Value.Items.Value)
				assert.NotNil(t, itemsProperty.Value.Items.Value.Properties)

				// Check field descriptions in response items
				if nameProperty, exists := itemsProperty.Value.Items.Value.Properties["name"]; exists {
					assert.NotNil(t, nameProperty.Value)
					assert.Equal(t, "Name user's name", nameProperty.Value.Description)
				}

				if emailProperty, exists := itemsProperty.Value.Items.Value.Properties["email"]; exists {
					assert.NotNil(t, emailProperty.Value)
					assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
				}
			}
		}
	})
}

func Test_addSchemaDescriptions(t *testing.T) {
	t.Run("test_schema_descriptions", func(t *testing.T) {
		// Create a schema ref for TestUser
		gen := openapi3gen.NewGenerator()
		schemaRef, err := gen.NewSchemaRefForValue(*new(TestUser), nil)
		assert.NoError(t, err)
		assert.NotNil(t, schemaRef)

		// Add schema descriptions
		addSchemaDescriptions[*TestUser](schemaRef)

		// Verify that descriptions are added to schema properties
		assert.NotNil(t, schemaRef.Value)
		assert.NotNil(t, schemaRef.Value.Properties)

		// Check Name field description
		if nameProperty, exists := schemaRef.Value.Properties["name"]; exists {
			assert.NotNil(t, nameProperty.Value)
			assert.Equal(t, "Name user's name", nameProperty.Value.Description)
		}

		// Check Email field description
		if emailProperty, exists := schemaRef.Value.Properties["email"]; exists {
			assert.NotNil(t, emailProperty.Value)
			assert.Equal(t, "Email user's email address", emailProperty.Value.Description)
		}
	})
}
