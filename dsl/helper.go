package dsl

import (
	"slices"
	"strings"

	"github.com/forbearing/golib/types/consts"
)

// rangeAction iterates through all actions in a Design and calls the provided function
// for each enabled action. This is a helper function used by Design.Range().
//
// The function iterates through actions in a predefined order and only processes
// actions that are enabled (action.Enabled == true).
//
// Parameters:
//   - d: The Design containing actions to iterate over
//   - fn: Callback function that receives (endpoint, action, phase) for each enabled action
//
// Iteration order:
//  1. Single record operations: Create, Delete, Update, Patch, List, Get
//  2. Batch operations: CreateMany, DeleteMany, UpdateMany, PatchMany
//  3. Data transfer operations: Import, Export
//
// For each enabled action, the callback receives:
//   - endpoint: The API endpoint path from the Design
//   - action: The Action configuration
//   - phase: The corresponding consts.Phase for the operation type
//
// Example:
//
//	rangeAction(design, func(ep string, a *Action, p consts.Phase) {
//		fmt.Printf("%s %s payload=%s result=%s\n", p.MethodName(), ep, a.Payload, a.Result)
//	})
func rangeAction(d *Design, fn func(string, *Action, consts.Phase)) {
	if d == nil || fn == nil || !d.Enabled {
		return
	}

	if d.Create.Enabled {
		fn(d.Endpoint, d.Create, consts.PHASE_CREATE)
	}
	if d.Delete.Enabled {
		fn(d.Endpoint, d.Delete, consts.PHASE_DELETE)
	}
	if d.Update.Enabled {
		fn(d.Endpoint, d.Update, consts.PHASE_UPDATE)
	}
	if d.Patch.Enabled {
		fn(d.Endpoint, d.Patch, consts.PHASE_PATCH)
	}
	if d.List.Enabled {
		fn(d.Endpoint, d.List, consts.PHASE_LIST)
	}
	if d.Get.Enabled {
		fn(d.Endpoint, d.Get, consts.PHASE_GET)
	}
	if d.CreateMany.Enabled {
		fn(d.Endpoint, d.CreateMany, consts.PHASE_CREATE_MANY)
	}
	if d.DeleteMany.Enabled {
		fn(d.Endpoint, d.DeleteMany, consts.PHASE_DELETE_MANY)
	}
	if d.UpdateMany.Enabled {
		fn(d.Endpoint, d.UpdateMany, consts.PHASE_UPDATE_MANY)
	}
	if d.PatchMany.Enabled {
		fn(d.Endpoint, d.PatchMany, consts.PHASE_PATCH_MANY)
	}
	if d.Import.Enabled {
		fn(d.Endpoint, d.Import, consts.PHASE_IMPORT)
	}
	if d.Export.Enabled {
		fn(d.Endpoint, d.Export, consts.PHASE_EXPORT)
	}
}

// is checks if the given name is a valid DSL method name.
// It verifies against the predefined list of supported DSL methods.
//
// Parameters:
//   - name: The method name to check
//
// Returns:
//   - bool: true if the name is a valid DSL method, false otherwise
func is(name string) bool {
	return slices.Contains(methodList, name)
}

// trimQuote removes surrounding quotes from a string.
// It trims double quotes ("), single quotes ('), and backticks (`) from both ends.
//
// Parameters:
//   - str: The string to trim quotes from
//
// Returns:
//   - string: The string with surrounding quotes removed
//
// Example:
//
//	trimQuote(`"hello"`) returns "hello"
//	trimQuote("'world'") returns "world"
func trimQuote(str string) string {
	return strings.TrimFunc(str, func(r rune) bool {
		return r == '`' || r == '"' || r == '\''
	})
}
