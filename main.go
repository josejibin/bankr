package main

import (
	"fmt"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	flag "github.com/spf13/pflag"

	log "github.com/Sirupsen/logrus"
)

var (

	// Global configuration reader.
	ko = koanf.New(".")

	// Version of the build.
	// This is injected at build-time.
	// Be sure to run the provided run script to inject correctly.
	buildVersion = "unknown"
	buildDate    = "unknown"
)

func init() {
	// Initialize the app configuration
	initConfig()

	// Initialize logger
	initLogger()
}

// Initializes the app configuration
func initConfig() {
	// Command line flags
	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Println(flagSet.FlagUsages())
		os.Exit(0)
	}

	// Commandline flags
	flagSet.String("config", "config.toml", "Path to the TOML configuration file")
	flagSet.Bool("version", false, "Current version of the build")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	// Load commandline params.
	ko.Load(posflag.Provider(flagSet, ".", ko), nil)

	// Display version.
	if ko.Bool("version") {
		fmt.Println(fmt.Sprintf("Commit: %v\nBuild: %v", buildVersion, buildDate))
		os.Exit(0)
	}

	// Load the config file.
	log.Printf("reading config: %s", ko.String("config"))
	if err := ko.Load(file.Provider(ko.String("config")), toml.Parser()); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

}

// Initialize loggers
func initLogger() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true})

	// Set log level based on environment
	if ko.Bool("app.debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	log.Debug("Current env : ", ko.Bool("app.debug"))

	// Initialize search
	initSearch()

	// Initialize server
	initServer(ko.String("server.address"))
}
