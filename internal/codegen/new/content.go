package new

import (
	"github.com/forbearing/golib/types/consts"
)

var modelContent = consts.CodeGeneratedComment() + `
package model

func init() {
}
`

var serviceContent = consts.CodeGeneratedComment() + `
package service

func init() {
}
`

var routerContent = consts.CodeGeneratedComment() + `
package router

func Init() error {
	return nil
}
`

var mainContent = consts.CodeGeneratedComment() + `
package main

import (
	"%s/configx"
	"%s/cronjobx"
	"%s/middlewarex"
	_ "%s/model"
	"%s/router"
	_ "%s/service"

	"github.com/forbearing/golib/bootstrap"
	. "github.com/forbearing/golib/util"
)

func main() {
	RunOrDie(middlewarex.Init)
	RunOrDie(bootstrap.Bootstrap)
	RunOrDie(configx.Init)
	RunOrDie(cronjobx.Init)
	RunOrDie(router.Init)
	RunOrDie(bootstrap.Run)
}
`

const configxContent = `package configx

func Init() error {
	return nil
}
`

const cronjobxContent = `package cronjobx

func Init() error {
	return nil
}
`

const middlewarexContent = `package middlewarex

func Init() error {
	return nil
}
`

const gitignoreContent = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Log files
*.log

# Temporary files
tmp/
temp/

# Build output
dist/
build/`
