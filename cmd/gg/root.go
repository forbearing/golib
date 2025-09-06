package main

import "github.com/spf13/cobra"

var (
	modelDir   string
	serviceDir string
	routerDir  string
	daoDir     string
	excludes   []string
	module     string
	debug      bool
	prune      bool
)

var rootCmd = &cobra.Command{
	Use:     "gg",
	Short:   "golib code generator",
	Long:    "golib code generator",
	Version: "1.0.0",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&module, "module", "M", "", "go project module name")
	rootCmd.PersistentFlags().StringVarP(&modelDir, "model", "m", "model", "model directory path")
	rootCmd.PersistentFlags().StringVarP(&serviceDir, "service", "s", "service", "service directory path")
	rootCmd.PersistentFlags().StringVarP(&routerDir, "router", "r", "router", "router directory path")
	rootCmd.PersistentFlags().StringVarP(&daoDir, "dao", "", "dao", "dao directory path")
	rootCmd.PersistentFlags().StringSliceVarP(&excludes, "exclude", "e", nil, "exclude files or directories")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&prune, "prune", false, "Prune disabled service action files with user confirmation")

	rootCmd.AddCommand(genCmd, newCmd, astCmd, pruneCmd, checkCmd, routesCmd, dockerCmd, k8sCmd, buildCmd, releaseCmd)
}
