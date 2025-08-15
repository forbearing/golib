package config

import (
	"runtime"
	"runtime/debug"
	"time"

	"github.com/forbearing/golib/types/consts"
)

const (
	// App related environment variables
	APP_NAME        = "APP_NAME"
	APP_VERSION     = "APP_VERSION"
	APP_DESCRIPTION = "APP_DESCRIPTION"
	APP_AUTHOR      = "APP_AUTHOR"
	APP_EMAIL       = "APP_EMAIL"
	APP_HOMEPAGE    = "APP_HOMEPAGE"
	APP_LICENSE     = "APP_LICENSE"
	APP_BUILD_TIME  = "APP_BUILD_TIME"
	APP_GIT_COMMIT  = "APP_GIT_COMMIT"
	APP_GIT_BRANCH  = "APP_GIT_BRANCH"
	APP_GO_VERSION  = "APP_GO_VERSION"
)

// AppInfo represents application metadata and build information
// This struct contains essential project information that can be used
// for version reporting, monitoring, and application identification
type AppInfo struct {
	// Basic application information
	Name        string `json:"name" mapstructure:"name" ini:"name" yaml:"name"`
	Version     string `json:"version" mapstructure:"version" ini:"version" yaml:"version"`
	Description string `json:"description" mapstructure:"description" ini:"description" yaml:"description"`
	Author      string `json:"author" mapstructure:"author" ini:"author" yaml:"author"`
	Email       string `json:"email" mapstructure:"email" ini:"email" yaml:"email"`
	Homepage    string `json:"homepage" mapstructure:"homepage" ini:"homepage" yaml:"homepage"`
	License     string `json:"license" mapstructure:"license" ini:"license" yaml:"license"`

	// Build and runtime information
	BuildTime  time.Time `json:"build_time" mapstructure:"build_time" ini:"build_time" yaml:"build_time"`
	GitCommit  string    `json:"git_commit" mapstructure:"git_commit" ini:"git_commit" yaml:"git_commit"`
	GitBranch  string    `json:"git_branch" mapstructure:"git_branch" ini:"git_branch" yaml:"git_branch"`
	GoVersion  string    `json:"go_version" mapstructure:"go_version" ini:"go_version" yaml:"go_version"`
	GitTag     string    `json:"git_tag" mapstructure:"git_tag" ini:"git_tag" yaml:"git_tag"`
	DirtyBuild bool      `json:"dirty_build" mapstructure:"dirty_build" ini:"dirty_build" yaml:"dirty_build"`
}

// setDefault sets default values for AppInfo configuration
func (a *AppInfo) setDefault() {
	if a.Name == "" {
		a.Name = consts.FrameworkName
	}
	if a.Version == "" {
		a.Version = "dev"
	}
	if a.Description == "" {
		a.Description = "A Go application built with golib framework"
	}
	if a.License == "" {
		a.License = "MIT"
	}
	if a.GoVersion == "" {
		a.GoVersion = runtime.Version()
	}

	// Try to get build info from runtime
	a.setBuildInfo()
}

// setBuildInfo attempts to extract build information from runtime/debug
func (a *AppInfo) setBuildInfo() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// Extract version control information from build settings
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			if a.GitCommit == "" {
				a.GitCommit = setting.Value
			}
		case "vcs.time":
			if a.BuildTime.IsZero() {
				if t, err := time.Parse(time.RFC3339, setting.Value); err == nil {
					a.BuildTime = t
				}
			}
		case "vcs.modified":
			a.DirtyBuild = setting.Value == "true"
		}
	}

	// Use module version if available and no custom version is set
	if a.Version == "dev" && buildInfo.Main.Version != "(devel)" && buildInfo.Main.Version != "" {
		a.Version = buildInfo.Main.Version
		a.GitTag = buildInfo.Main.Version
	}
}

// GetShortCommit returns the short version of git commit hash
func (a *AppInfo) GetShortCommit() string {
	if len(a.GitCommit) > 7 {
		return a.GitCommit[:7]
	}
	return a.GitCommit
}

// GetVersionInfo returns formatted version information
func (a *AppInfo) GetVersionInfo() string {
	version := a.Version
	if a.GitCommit != "" {
		version += "+" + a.GetShortCommit()
	}
	if a.DirtyBuild {
		version += "-dirty"
	}
	return version
}

// GetBuildInfo returns comprehensive build information as a map
func (a *AppInfo) GetBuildInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":        a.Name,
		"version":     a.Version,
		"description": a.Description,
		"author":      a.Author,
		"email":       a.Email,
		"homepage":    a.Homepage,
		"license":     a.License,
		"build_time":  a.BuildTime,
		"git_commit":  a.GitCommit,
		"git_branch":  a.GitBranch,
		"git_tag":     a.GitTag,
		"go_version":  a.GoVersion,
		"dirty_build": a.DirtyBuild,
	}
}
