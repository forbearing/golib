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
