package bootstrap

import (
	"fmt"
	"os"
)

var handlers = []func(){}

// Exit will call all registered cleanup handlers and then exit.
func Exit(code int) {
	runHandlers()
	os.Exit(code)
}

// RegisterCleanup append custom cleanup handler, the handler will be invoked by `Cleanup` function.
// first handler will be called first
func RegisterCleanup(handler func()) {
	handlers = append(handlers, handler)
}

// DeferCleanup same as RegisterCleanup, but last handler will be called first.
func DeferCleanup(handler func()) {
	handlers = append([]func(){handler}, handlers...)
}

// Cleanup will call all registered cleanup handlers.
func Cleanup() {
	runHandlers()
}

func runHandlers() {
	for _, handler := range handlers {
		go safeRun(handler)
	}
}

func safeRun(handler func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: cleanup handler error:", err)
		}
	}()

	handler()
}
