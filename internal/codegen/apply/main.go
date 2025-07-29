package apply

import (
	"fmt"
	"os"

	logrusLib "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	modelDir   string
	serviceDir string
	excludes   []string
	module     string
	debug      bool
)

func init() {
	pflag.StringVarP(&modelDir, "model", "m", "model", "model directory path")
	pflag.StringVarP(&serviceDir, "service", "s", "service", "service directory path")
	pflag.StringVarP(&module, "module", "M", "", "module path")
	pflag.StringSliceVarP(&excludes, "exclude", "e", nil, "exclude files")
	pflag.BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	pflag.Parse()

	// Initialize default logger (will be replaced in Main if needed)
	initDefaultLogger()
}

// initDefaultLogger initializes a basic logger for early use
func initDefaultLogger() {
	// Set default log level to Info
	logrusLib.SetLevel(logrusLib.InfoLevel)
	logrusLib.SetOutput(os.Stdout)
	logrusLib.SetFormatter(&logrusLib.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger = logrusLib.WithField("component", "codegen-apply")
	parserLogger = logrusLib.WithField("component", "codegen-apply-parser")
}

// InitLogger initializes the console logger based on debug flag
func InitLogger(debug bool) {
	var level logrusLib.Level
	if debug {
		level = logrusLib.DebugLevel
	} else {
		level = logrusLib.InfoLevel
	}

	logrusLib.SetLevel(level)
	logrusLib.SetOutput(os.Stdout)
	logrusLib.SetFormatter(&logrusLib.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Replace the global loggers
	logger = logrusLib.WithField("component", "codegen-apply")
	parserLogger = logrusLib.WithField("component", "codegen-apply-parser")
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
