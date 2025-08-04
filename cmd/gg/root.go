package main

import "github.com/spf13/cobra"

var (
	modelDir   string
	serviceDir string
	excludes   []string
	module     string
	debug      bool
)

var rootCmd = &cobra.Command{
	Use:     "gg",
	Short:   "golib code generator",
	Long:    "golib code generator",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		genRun()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&module, "module", "M", "", "go project module name")
	rootCmd.PersistentFlags().StringVarP(&modelDir, "model", "m", "model", "model directory path")
	rootCmd.PersistentFlags().StringVarP(&serviceDir, "service", "s", "service", "service directory path")
	rootCmd.PersistentFlags().StringSliceVarP(&excludes, "exclude", "e", nil, "exclude files or directories")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	rootCmd.AddCommand(genCmd, applyCmd, newCmd)
}
