package dsl

import "github.com/forbearing/golib/types/consts"

// RangeAction walks through all enabled actions in the given Design and invokes
// the provided callback fn for each, in a fixed order.
//
// Behavior:
//   - If design is nil, fn is nil, or design.Enabled is false, the function returns immediately.
//   - Only actions with Enabled == true are passed to the callback.
//   - Invocation order is fixed: Create, Delete, Update, Patch, List, Get,
//     CreateMany, DeleteMany, UpdateMany, PatchMany.
//   - For each invocation, the callback receives:
//     endpoint string     -> design.Endpoint
//     act *Action         -> the current action
//     phase consts.Phase  -> the phase corresponding to the current action
//   - Execution is sequential; the function does not call the callback concurrently.
//
// Example:
//
//	RangeAction(design, func(ep string, a *Action, p consts.Phase) {
//	    fmt.Printf("%s %s payload=%s result=%s\n", p.MethodName(), ep, a.Payload, a.Result)
//	})
func RangeAction(design *Design, fn func(string, *Action, consts.Phase)) {
	if design == nil || fn == nil {
		return
	}

	if !design.Enabled {
		return
	}

	if design.Create.Enabled {
		fn(design.Endpoint, design.Create, consts.PHASE_CREATE)
	}
	if design.Delete.Enabled {
		fn(design.Endpoint, design.Delete, consts.PHASE_DELETE)
	}
	if design.Update.Enabled {
		fn(design.Endpoint, design.Update, consts.PHASE_UPDATE)
	}
	if design.Patch.Enabled {
		fn(design.Endpoint, design.Patch, consts.PHASE_PATCH)
	}
	if design.List.Enabled {
		fn(design.Endpoint, design.List, consts.PHASE_LIST)
	}
	if design.Get.Enabled {
		fn(design.Endpoint, design.Get, consts.PHASE_GET)
	}
	if design.CreateMany.Enabled {
		fn(design.Endpoint, design.CreateMany, consts.PHASE_CREATE_MANY)
	}
	if design.DeleteMany.Enabled {
		fn(design.Endpoint, design.DeleteMany, consts.PHASE_DELETE_MANY)
	}
	if design.UpdateMany.Enabled {
		fn(design.Endpoint, design.UpdateMany, consts.PHASE_UPDATE_MANY)
	}
	if design.PatchMany.Enabled {
		fn(design.Endpoint, design.PatchMany, consts.PHASE_PATCH_MANY)
	}
	if design.Import.Enabled {
		fn(design.Endpoint, design.Import, consts.PHASE_IMPORT)
	}
	if design.Export.Enabled {
		fn(design.Endpoint, design.Export, consts.PHASE_EXPORT)
	}
}
