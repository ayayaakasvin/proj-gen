package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/ayayaakasvin/proj-gen/internal/lib/prompt"
	"github.com/ayayaakasvin/proj-gen/internal/options"
	"github.com/ayayaakasvin/proj-gen/templates"
	"github.com/sirupsen/logrus"
)

var overwrite bool

type Generator struct {
	opts   *options.ProjectOptions
	logger *logrus.Logger
}

type GenerateFunc func() error

func GenerateProject(opts *options.ProjectOptions, logger *logrus.Logger) {
	g := &Generator{
		opts:   opts,
		logger: logger,
	}

	err := g.Generate()
	if err != nil {
		FallBack(path.Join(opts.Path, opts.Name))
	}
}

func (g *Generator) Generate() error {
	funcSlice := []GenerateFunc{
		g.GenerateGoMod,
		g.GenerateMain,
		g.GenerateConfig,
		g.RunGoModTidy,
	}

	for _, f := range funcSlice {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) GenerateGoMod() error {
	g.logger.Info("Generating go.mod")

	pathToProject := path.Join(g.opts.Path, g.opts.Name)

	if _, err := os.Stat(pathToProject); !os.IsNotExist(err) {
		g.logger.Warn("Project already exists")
		if !prompt.AskForOverwrite() {
			g.logger.Info("Exiting application")
			os.Exit(0)
		} else {
			g.logger.Info("Overwriting project")
			overwrite = true
		}
	}

	g.logger.Info("Creating project directory")
	err := os.MkdirAll(pathToProject, 0755)
	if err != nil {
		g.logger.Error("Error creating project directory: ", err)
		return err
	}

	g.logger.Info("Creating go.mod file")
	goModFile, err := os.Create(path.Join(pathToProject, "go.mod"))
	if err != nil {
		g.logger.Error("Error creating go.mod file: ", err)
		return err
	}
	defer goModFile.Close()

	tmpl, err := templates.GetTemplate("go.mod.tmpl")
	if err != nil {
		g.logger.Error("Error reading template file: ", err)
		return err
	}

	err = tmpl.Execute(goModFile, g.opts)
	if err != nil {
		g.logger.Error("Error executing template: ", err)
		return err
	}

	return nil
}

func (g *Generator) GenerateMain() error {
	g.logger.Info("Generating main.go")

	pathToProject := path.Join(g.opts.Path, g.opts.Name)
	mainFile, err := os.Create(path.Join(pathToProject, "main.go"))
	if err != nil {
		g.logger.Error("Error creating main.go file: ", err)
		return err
	}
	defer mainFile.Close()

	tmpl, err := templates.GetTemplate("main.go.tmpl")
	if err != nil {
		g.logger.Error("Error reading template file: ", err)
		return err
	}

	err = tmpl.Execute(mainFile, nil)
	if err != nil {
		g.logger.Error("Error executing template: ", err)
		return err
	}

	return nil
}

func (g *Generator) GenerateConfig() error {
	if !g.opts.WithConfig {
		return nil
	}

	g.logger.Info("Generating config directory")

	pathToProject := path.Join(g.opts.Path, g.opts.Name)
	configDir := path.Join(pathToProject, "config")

	g.logger.Info("Creating config directory")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		g.logger.Error("Error creating config directory: ", err)
		return err
	}

	g.logger.Info("Creating config.yaml file")
	configFile, err := os.Create(path.Join(configDir, "config.yaml"))
	if err != nil {
		g.logger.Error("Error creating config.yaml file: ", err)
		return err
	}
	defer configFile.Close()

	g.logger.Info("Creating internal/config directory")
	configInternalDir := path.Join(pathToProject, "internal", "config")
	err = os.MkdirAll(configInternalDir, 0755)
	if err != nil {
		g.logger.Error("Error creating internal/config directory: ", err)
		return err
	}
	g.logger.Info("Creating config.go file")
	configGoFile, err := os.Create(path.Join(configInternalDir, "config.go"))
	if err != nil {
		g.logger.Error("Error creating config.go file: ", err)
		return err
	}
	defer configGoFile.Close()

	tmpl, err := templates.GetTemplate("config.go.tmpl")
	if err != nil {
		g.logger.Error("Error reading template file: ", err)
		return err
	}

	err = tmpl.Execute(configGoFile, nil)
	if err != nil {
		g.logger.Error("Error executing template: ", err)
		return err
	}

	return nil
}

func FallBack(pathToProject string) {
	if !overwrite {
		err := os.RemoveAll(pathToProject)
		if err != nil {
			fmt.Println("Error removing project directory: ", err)
			return
		}
	}

	os.Exit(1)
}

func (g *Generator) RunGoModTidy() error {
	if !g.opts.RunTidy {
		return nil
	}

	g.logger.Info("Running go mod tidy")

	pathToProject := path.Join(g.opts.Path, g.opts.Name)

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = pathToProject

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running go mod tidy: ", err)
		return err
	}

	return nil
}
