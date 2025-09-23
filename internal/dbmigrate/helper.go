package dbmigrate

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	// Color tools - consistent with cmd/gg
	green   = color.New(color.FgHiGreen).SprintFunc()
	yellow  = color.New(color.FgHiYellow).SprintFunc()
	red     = color.New(color.FgHiRed).SprintFunc()
	cyan    = color.New(color.FgHiCyan).SprintFunc()
	magenta = color.New(color.FgHiMagenta).SprintFunc()
	blue    = color.New(color.FgHiBlue).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
)

// logSection outputs section title
func logSection(title string) {
	fmt.Printf("%s %s\n", cyan("▶"), bold(title))
}

// logSuccess outputs success message
func logSuccess(msg string) {
	fmt.Printf("  %s %s\n", green("✔"), msg)
}

// logError outputs error message
func logError(msg string) {
	fmt.Printf("  %s %s\n", red("✘"), msg)
}

// logWarning outputs warning message
func logWarning(msg string) {
	fmt.Printf("  %s %s\n", yellow("⚠"), msg)
}

// logSeparator outputs separator line
func logSeparator(char string, length int) {
	fmt.Println(strings.Repeat(char, length))
}

// logTitle outputs title with separator
func logTitle(title string, count int) {
	fmt.Printf("%s (%d changes):\n", title, count)
}

// logOperation outputs operation with index
func logOperation(index int, cmd string) {
	fmt.Printf("%d. %s\n", index, cmd)
}

// logPrompt outputs prompt without newline
func logPrompt(prompt string) {
	fmt.Printf("%s (y/N): ", prompt)
}
