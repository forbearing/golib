package main

import (
	"github.com/forbearing/golib/internal/codegen/apply"
	"github.com/forbearing/golib/internal/codegen/gen"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply service code",
	Run: func(cmd *cobra.Command, args []string) {
		if len(module) == 0 {
			var err error
			module, err = gen.GetModulePath()
			checkErr(err)
		}
		config := apply.NewApplyConfig(module, modelDir, serviceDir).WithExclusions(excludes...)

		checkErr(apply.ApplyServiceGeneration(config))
	},
}
