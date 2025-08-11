package consts_test

import (
	"testing"

	"github.com/forbearing/golib/types/consts"
)

func TestPhase_MethodName(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		phase consts.Phase
		want  string
	}{
		{
			name:  "CreateBefore",
			phase: consts.PHASE_CREATE_BEFORE,
			want:  "CreateBefore",
		},
		{
			name:  "UpdateManyAfter",
			phase: consts.PHASE_UPDATE_MANY_AFTER,
			want:  "UpdateManyAfter",
		},
		{
			name:  "PatchMany",
			phase: consts.PHASE_PATCH_MANY,
			want:  "PatchMany",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.phase.MethodName()
			if got != tt.want {
				t.Errorf("MethodName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhase_RoleName(t *testing.T) {
	tests := []struct {
		name  string
		phase consts.Phase
		want  string
	}{
		// 单个 CRUD
		{"creator", consts.PHASE_CREATE, "Creator"},
		{"updater", consts.PHASE_UPDATE, "Updater"},
		{"deleter", consts.PHASE_DELETE, "Deleter"},
		{"patcher", consts.PHASE_PATCH, "Patcher"},
		{"lister", consts.PHASE_LIST, "Lister"},
		{"getter", consts.PHASE_GET, "Getter"},

		// Many CRUD
		{"many_creator", consts.PHASE_CREATE_MANY, "ManyCreator"},
		{"many_updater", consts.PHASE_UPDATE_MANY, "ManyUpdater"},
		{"many_deleter", consts.PHASE_DELETE_MANY, "ManyDeleter"},
		{"many_patcher", consts.PHASE_PATCH_MANY, "ManyPatcher"},

		// before/after - 单个
		{"create_before", consts.PHASE_CREATE_BEFORE, "Creator"},
		{"create_after", consts.PHASE_CREATE_AFTER, "Creator"},
		{"update_before", consts.PHASE_UPDATE_BEFORE, "Updater"},
		{"update_after", consts.PHASE_UPDATE_AFTER, "Updater"},
		{"delete_before", consts.PHASE_DELETE_BEFORE, "Deleter"},
		{"delete_after", consts.PHASE_DELETE_AFTER, "Deleter"},
		{"patch_before", consts.PHASE_PATCH_BEFORE, "Patcher"},
		{"patch_after", consts.PHASE_PATCH_AFTER, "Patcher"},
		{"list_before", consts.PHASE_LIST_BEFORE, "Lister"},
		{"list_after", consts.PHASE_LIST_AFTER, "Lister"},
		{"get_before", consts.PHASE_GET_BEFORE, "Getter"},
		{"get_after", consts.PHASE_GET_AFTER, "Getter"},

		// before/after - Many
		{"many_create_before", consts.PHASE_CREATE_MANY_BEFORE, "ManyCreator"},
		{"many_create_after", consts.PHASE_CREATE_MANY_AFTER, "ManyCreator"},
		{"many_update_before", consts.PHASE_UPDATE_MANY_BEFORE, "ManyUpdater"},
		{"many_update_after", consts.PHASE_UPDATE_MANY_AFTER, "ManyUpdater"},
		{"many_delete_before", consts.PHASE_DELETE_MANY_BEFORE, "ManyDeleter"},
		{"many_delete_after", consts.PHASE_DELETE_MANY_AFTER, "ManyDeleter"},
		{"many_patch_before", consts.PHASE_PATCH_MANY_BEFORE, "ManyPatcher"},
		{"many_patch_after", consts.PHASE_PATCH_MANY_AFTER, "ManyPatcher"},

		// 非 CRUD 操作
		{"filter", consts.PHASE_FILTER, ""},
		{"import", consts.PHASE_IMPORT, ""},
		{"export", consts.PHASE_EXPORT, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.phase.RoleName()
			if got != tt.want {
				t.Errorf("RoleName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhase_BeforeAfter(t *testing.T) {
	tests := []struct {
		name       string
		phase      consts.Phase
		wantBefore consts.Phase
		wantAfter  consts.Phase
	}{
		// Single CRUD
		{"create", consts.PHASE_CREATE, consts.PHASE_CREATE_BEFORE, consts.PHASE_CREATE_AFTER},
		{"update", consts.PHASE_UPDATE, consts.PHASE_UPDATE_BEFORE, consts.PHASE_UPDATE_AFTER},
		{"delete", consts.PHASE_DELETE, consts.PHASE_DELETE_BEFORE, consts.PHASE_DELETE_AFTER},
		{"patch", consts.PHASE_PATCH, consts.PHASE_PATCH_BEFORE, consts.PHASE_PATCH_AFTER},
		{"list", consts.PHASE_LIST, consts.PHASE_LIST_BEFORE, consts.PHASE_LIST_AFTER},
		{"get", consts.PHASE_GET, consts.PHASE_GET_BEFORE, consts.PHASE_GET_AFTER},

		// Many CRUD
		{"create_many", consts.PHASE_CREATE_MANY, consts.PHASE_CREATE_MANY_BEFORE, consts.PHASE_CREATE_MANY_AFTER},
		{"update_many", consts.PHASE_UPDATE_MANY, consts.PHASE_UPDATE_MANY_BEFORE, consts.PHASE_UPDATE_MANY_AFTER},
		{"delete_many", consts.PHASE_DELETE_MANY, consts.PHASE_DELETE_MANY_BEFORE, consts.PHASE_DELETE_MANY_AFTER},
		{"patch_many", consts.PHASE_PATCH_MANY, consts.PHASE_PATCH_MANY_BEFORE, consts.PHASE_PATCH_MANY_AFTER},

		// Already before/after → no change
		{"create_before", consts.PHASE_CREATE_BEFORE, consts.PHASE_CREATE_BEFORE, consts.PHASE_CREATE_BEFORE},
		{"update_after", consts.PHASE_UPDATE_AFTER, consts.PHASE_UPDATE_AFTER, consts.PHASE_UPDATE_AFTER},
		{"many_delete_before", consts.PHASE_DELETE_MANY_BEFORE, consts.PHASE_DELETE_MANY_BEFORE, consts.PHASE_DELETE_MANY_BEFORE},
		{"many_patch_after", consts.PHASE_PATCH_MANY_AFTER, consts.PHASE_PATCH_MANY_AFTER, consts.PHASE_PATCH_MANY_AFTER},

		// Non CRUD → no change
		{"filter", consts.PHASE_FILTER, consts.PHASE_FILTER, consts.PHASE_FILTER},
		{"import", consts.PHASE_IMPORT, consts.PHASE_IMPORT, consts.PHASE_IMPORT},
		{"export", consts.PHASE_EXPORT, consts.PHASE_EXPORT, consts.PHASE_EXPORT},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.phase.Before(); got != tt.wantBefore {
				t.Errorf("Before() = %v, want %v", got, tt.wantBefore)
			}
			if got := tt.phase.After(); got != tt.wantAfter {
				t.Errorf("After() = %v, want %v", got, tt.wantAfter)
			}
		})
	}
}
