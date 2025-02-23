package bootstrap

import (
	"fmt"
	"os"
)

var handlers = []func(){}

func runHandler(handler func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: exit handler error:", err)
		}
	}()

	handler()
}

func runHandlers() {
	for _, handler := range handlers {
		runHandler(handler)
	}
}

// Exit will call all exit handlers and then exit.
func Exit(code int) {
	runHandlers()
	os.Exit(code)
}

// RegisterExitHandler append custom exit handler, the handler will be invoked by `Exit(int)` function.
// first handler will be called first
func RegisterExitHandler(handler func()) {
	handlers = append(handlers, handler)
}

// DeferExitHandler same as RegisterExitHandler, but last handler will be called first.
func DeferExitHandler(handler func()) {
	handlers = append([]func(){handler}, handlers...)
}

// Cleanup
func Cleanup() {
	runHandlers()
}
