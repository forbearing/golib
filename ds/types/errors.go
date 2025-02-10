package types

import (
	"fmt"
)

var (
	ErrComparisonNil = fmt.Errorf("comparison function cannot be nil")
	ErrEqualNil      = fmt.Errorf("equality function cannot be nil")
	ErrFuncNil       = fmt.Errorf("function cannot be nil")
)
