package options

import (
	"flag"
	"fmt"
	"os"

	goversion "github.com/ayayaakasvin/gen-project/internal/lib/go-version"
)

type ProjectOptions struct {
	Name 		string
	Module 		string
	WithConfig	bool
	RunTidy 	bool
	GoVersion	string
	Path 		string
}

// AfterParseFunc is a function type that takes a pointer to ProjectOptions
type AfterParseFunc func(*ProjectOptions)

func ParseFlags() (*ProjectOptions, AfterParseFunc) {
	name := flag.String("name", "", "Project name (required)")
	module := flag.String("module", "", "Project module name (default: project name)")
	
	var withConfig *bool = new(bool)
	flag.BoolVar(withConfig, "config", false, "Include config directory")
	flag.BoolVar(withConfig, "c", false, "Alias for --config")
	
	var runTidy *bool = new(bool)
	flag.BoolVar(runTidy, "tidy", false, "Run go mod tidy")
	flag.BoolVar(runTidy, "t", false, "Alias for --tidy")
	
	goVersion := flag.String("go-version", "", "Go version to use(default: build version)")
	path := flag.String("path", "./", "Path to create the project (default: current directory)")

	opts := &ProjectOptions{}

	afterParse := func(o *ProjectOptions) {
		o.Name = *name
		o.Module = *module
		o.WithConfig = *withConfig
		o.GoVersion = *goVersion
		o.Path = *path
		o.RunTidy = *runTidy

		if o.Name == "" {
			flag.Usage()
			fmt.Fprintln(os.Stderr, "Error: -name is required")
			os.Exit(1)
		}
		if o.Module == "" {
			o.Module = o.Name
		}
		if o.GoVersion == "" {
			o.GoVersion = goversion.Version()
		}
		
	}

	return opts, afterParse
}

func (o *ProjectOptions) String() string {
	return fmt.Sprintf("ProjectOptions{Name: %s, Module: %s, WithConfig: %t, Go-Version: %s, Path: %s}", 
		o.Name, 
		o.Module, 
		o.WithConfig, 
		o.GoVersion, 
		o.Path,
	)
}