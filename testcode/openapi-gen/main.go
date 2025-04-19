package main

import (
	"fmt"
	"log"

	"wokao/pkg/openapigen"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	fmt.Println("Generating OpenAPI spec...")
	var doc *openapi3.T
	var err error
	if doc, err = openapigen.GenerateOpenAPI(); err != nil {
		log.Fatalf("Failed to generate OpenAPI spec: %v", err)
	}

	// 设置 Gin 路由
	r := gin.Default()

	// 提供 Swagger JSON
	r.GET("/api.json", gin.WrapH(openapigen.DocumentHandler(doc)))

	// 添加 Swagger UI 路由
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api.json")))

	// 添加其他 API 路由
	// ...

	// 启动服务器
	fmt.Println("API documentation available at http://localhost:8003/docs/index.html")
	r.Run(":8003")
}
