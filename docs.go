// Package docs 提供 golib 框架的 API 文档
package docs

// @title           Golib API
// @version         v0.4.2
// @description     API documentation for the golib framework
// @termsOfService  http://swagger.io/terms/

// @contact.name   Golib Support
// @contact.url    https://github.com/forbearing/golib/issues
// @contact.email  your-email@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// 用户相关 API
type userAPI struct{}

// @Summary      创建用户
// @Description  创建新用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      UserInput  true  "User Info"
// @Success      201   {object}  User
// @Failure      400   {object}  ErrorResponse
// @Router       /users [post]
func (api *userAPI) CreateUser() {}

// @Summary      获取用户列表
// @Description  获取所有用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page     query   int  false  "Page number"      default(1)
// @Param        pageSize query   int  false  "Items per page"   default(10)
// @Success      200      {array}  User
// @Router       /users [get]
func (api *userAPI) ListUsers() {}

// @Summary      获取用户
// @Description  通过ID获取用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  User
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
func (api *userAPI) GetUser() {}

// @Summary      更新用户
// @Description  更新用户信息
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string    true  "User ID"
// @Param        user  body      UserInput  true  "User Info"
// @Success      200   {object}  User
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Router       /users/{id} [put]
func (api *userAPI) UpdateUser() {}

// @Summary      删除用户
// @Description  删除用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [delete]
func (api *userAPI) DeleteUser() {}

// User 用户模型
type User struct {
	ID        string `json:"id" example:"1"`
	Name      string `json:"name" example:"John Doe"`
	Email     string `json:"email" example:"john@example.com"`
	CreatedAt string `json:"createdAt" example:"2025-04-15T10:00:00Z"`
	UpdatedAt string `json:"updatedAt" example:"2025-04-15T10:00:00Z"`
}

// UserInput 用户输入模型
type UserInput struct {
	Name  string `json:"name" example:"John Doe" binding:"required"`
	Email string `json:"email" example:"john@example.com" binding:"required,email"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"User not found"`
}
