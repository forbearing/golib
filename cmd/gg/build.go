package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build cross-platform binaries with version information",
	Long: `Build cross-platform binaries with embedded version information.

This command will:
- Build binaries for multiple platforms and architectures
- Embed version information using -ldflags
- Support custom output directory
- Generate build metadata

Examples:
  gg build                           # Build for current platform
  gg build --all                     # Build for all supported platforms
  gg build --os linux --arch amd64   # Build for specific platform
  gg build --output ./dist           # Custom output directory
  gg build --version v1.0.0          # Custom version`,
	RunE: buildRun,
}

var (
	buildAll     bool
	buildOS      []string
	buildArch    []string
	buildOutput  string
	buildVersion string
	buildLdflags string
)

// Supported platforms and architectures
var (
	supportedOS   = []string{"linux", "darwin", "windows"}
	supportedArch = []string{"amd64", "arm64"}
)

func init() {
	buildCmd.Flags().BoolVar(&buildAll, "all", false, "Build for all supported platforms")
	buildCmd.Flags().StringSliceVar(&buildOS, "os", nil, "Target operating systems (linux,darwin,windows)")
	buildCmd.Flags().StringSliceVar(&buildArch, "arch", nil, "Target architectures (amd64,arm64)")
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "./dist", "Output directory for binaries")
	buildCmd.Flags().StringVarP(&buildVersion, "version", "v", "", "Version to embed (default: auto-detect from git)")
	buildCmd.Flags().StringVar(&buildLdflags, "ldflags", "", "Additional ldflags to pass to go build")
}

func buildRun(cmd *cobra.Command, args []string) error {
	logSection("Build Binaries")

	// Get build information
	buildInfo, err := getBuildInfo()
	if err != nil {
		return fmt.Errorf("failed to get build info: %w", err)
	}

	// Determine target platforms
	targets, err := getBuildTargets()
	if err != nil {
		return fmt.Errorf("failed to determine build targets: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(buildOutput, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build for each target
	successCount := 0
	for _, target := range targets {
		if err := buildForTarget(target, buildInfo); err != nil {
			fmt.Printf("%s Failed to build for %s: %v\n", red("✘"), target.String(), err)
			continue
		}
		fmt.Printf("%s Built successfully for %s\n", green("✔"), target.String())
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("all builds failed")
	}

	fmt.Printf("%s Build completed: %d/%d successful\n", green("✔"), successCount, len(targets))
	fmt.Printf("%s Output directory: %s\n", gray("→"), buildOutput)

	return nil
}

// BuildTarget represents a build target
type BuildTarget struct {
	OS   string
	Arch string
}

func (t BuildTarget) String() string {
	return fmt.Sprintf("%s/%s", t.OS, t.Arch)
}

// BuildInfo contains build metadata
type BuildInfo struct {
	Version   string
	Commit    string
	Branch    string
	BuildTime string
	GoVersion string
	Module    string
	Binary    string
}

// getBuildInfo collects build information
func getBuildInfo() (*BuildInfo, error) {
	info := &BuildInfo{
		BuildTime: time.Now().UTC().Format(time.RFC3339),
		GoVersion: runtime.Version(),
	}

	// Get module name
	moduleName, err := getModuleName()
	if err != nil {
		return nil, fmt.Errorf("failed to get module name: %w", err)
	}
	info.Module = moduleName
	info.Binary = filepath.Base(moduleName)

	// Get version
	if buildVersion != "" {
		info.Version = buildVersion
	} else {
		version, err := getGitVersion()
		if err != nil {
			fmt.Printf("%s Failed to get git version, using 'dev': %v\n", yellow("⚠"), err)
			info.Version = "dev"
		} else {
			info.Version = version
		}
	}

	// Get git commit
	commit, err := getGitCommit()
	if err != nil {
		fmt.Printf("%s Failed to get git commit: %v\n", yellow("⚠"), err)
		info.Commit = "unknown"
	} else {
		info.Commit = commit
	}

	// Get git branch
	branch, err := getGitBranch()
	if err != nil {
		fmt.Printf("%s Failed to get git branch: %v\n", yellow("⚠"), err)
		info.Branch = "unknown"
	} else {
		info.Branch = branch
	}

	return info, nil
}

// getBuildTargets determines which platforms to build for
func getBuildTargets() ([]BuildTarget, error) {
	var targets []BuildTarget

	if buildAll {
		// Build for all supported platforms
		for _, os := range supportedOS {
			for _, arch := range supportedArch {
				targets = append(targets, BuildTarget{OS: os, Arch: arch})
			}
		}
	} else if len(buildOS) > 0 || len(buildArch) > 0 {
		// Build for specified platforms
		osTargets := buildOS
		if len(osTargets) == 0 {
			osTargets = []string{runtime.GOOS}
		}

		archTargets := buildArch
		if len(archTargets) == 0 {
			archTargets = []string{runtime.GOARCH}
		}

		for _, os := range osTargets {
			for _, arch := range archTargets {
				targets = append(targets, BuildTarget{OS: os, Arch: arch})
			}
		}
	} else {
		// Build for current platform only
		targets = append(targets, BuildTarget{OS: runtime.GOOS, Arch: runtime.GOARCH})
	}

	return targets, nil
}

// buildForTarget builds the binary for a specific target
func buildForTarget(target BuildTarget, info *BuildInfo) error {
	// Generate binary name
	binaryName := info.Binary
	if target.OS == "windows" {
		binaryName += ".exe"
	}

	// Generate output path
	outputPath := filepath.Join(buildOutput, fmt.Sprintf("%s_%s_%s", info.Binary, target.OS, target.Arch), binaryName)

	// Create target directory
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Build ldflags
	ldflags := buildLdflags
	if ldflags != "" {
		ldflags += " "
	}
	ldflags += fmt.Sprintf("-X 'main.version=%s' -X 'main.commit=%s' -X 'main.branch=%s' -X 'main.buildTime=%s' -X 'main.goVersion=%s'",
		info.Version, info.Commit, info.Branch, info.BuildTime, info.GoVersion)

	// Prepare build command
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", outputPath, ".")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", target.OS),
		fmt.Sprintf("GOARCH=%s", target.Arch),
		"CGO_ENABLED=0",
	)

	// Execute build
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// getGitVersion gets version from git tags
func getGitVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getGitCommit gets current git commit hash
func getGitCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getGitBranch gets current git branch
func getGitBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
