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
	pathToProject string
}

type GenerateFunc func() error

func GenerateProject(opts *options.ProjectOptions, logger *logrus.Logger) {
	g := &Generator{
		opts:   opts,
		logger: logger,
		pathToProject: path.Join(opts.Path, opts.Name),
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
		g.GenerateInternal,
		g.GenerateConfig,
		g.RunGoModTidy,
		g.GenerateGitIgnore,
		g.RunGitInit,
		g.GenerateDockerfile,
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
	if _, err := os.Stat(g.pathToProject); !os.IsNotExist(err) {
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
	err := os.MkdirAll(g.pathToProject, 0755)
	if err != nil {
		g.logger.Error("Error creating project directory: ", err)
		return err
	}

	g.logger.Info("Creating go.mod file")
	goModFile, err := os.Create(path.Join(g.pathToProject, "go.mod"))
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

	mainFile, err := os.Create(path.Join(g.pathToProject, "main.go"))
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

func (g *Generator) GenerateInternal() error {
	g.logger.Info("Generating internal directory")
	internalDir := path.Join(g.pathToProject, "internal")
	err := os.MkdirAll(internalDir, 0755)
	if err != nil {
		g.logger.Error("Error creating internal directory: ", err)
		return err
	}

	return nil
}

func (g *Generator) GenerateConfig() error {
	if !g.opts.WithConfig {
		return nil
	}

	g.logger.Info("Generating config directory")

	
	configDir := path.Join(g.pathToProject, "config")

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
	configInternalDir := path.Join(g.pathToProject, "internal", "config")
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

func (g *Generator) GenerateGitIgnore() error {
	if !g.opts.RunGitInit {
		return nil
	}
	g.logger.Info("Generating .gitignore file")

	gitIgnoreFile, err := os.Create(path.Join(g.pathToProject, ".gitignore"))
	if err != nil {
		g.logger.Error("Error creating .gitignore file: ", err)
		return err
	}
	defer gitIgnoreFile.Close()

	_, err = gitIgnoreFile.WriteString("*\n")
	if err != nil {
		g.logger.Error("Error writing to .gitignore file: ", err)
		return err
	}

	return nil
}

func (g *Generator) RunGoModTidy() error {
	if !g.opts.RunTidy {
		return nil
	}

	g.logger.Info("Running go mod tidy")

	

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = g.pathToProject

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running go mod tidy: ", err)
		return err
	}

	return nil
}

func (g *Generator) RunGitInit() error {
	if !g.opts.RunGitInit {
		return nil
	}

	g.logger.Info("Running git init")

	cmd := exec.Command("git", "init")
	cmd.Dir = g.pathToProject

	err := cmd.Run()
	if err != nil {
		g.logger.Error("Error running git init: ", err)
		return err
	}

	return nil
}

func (g *Generator) GenerateDockerfile() error {
	if !g.opts.WithDocker {
		return nil
	}
	g.logger.Info("Generating Dockerfile")
	dockerFile, err := os.Create(path.Join(g.pathToProject, "Dockerfile"))
	if err != nil {
		g.logger.Error("Error creating Dockerfile: ", err)
		return err
	}
	defer dockerFile.Close()

	tmpl, err := templates.GetTemplate("Dockerfile.tmpl")
	if err != nil {
		g.logger.Error("Error reading template file: ", err)
		return err
	}

	err = tmpl.Execute(dockerFile, g.opts)
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