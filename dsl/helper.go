package dsl

import "github.com/forbearing/golib/types/consts"

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
}
