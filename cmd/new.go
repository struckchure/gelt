package main

import (
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/struckchure/gelt"
)

var qs = []*survey.Question{
	{
		Name:      "project_dir",
		Prompt:    &survey.Input{Message: "Project Directory:", Default: "."},
		Transform: survey.ToLower,
	},
	{
		Name: "package_name",
		Prompt: &survey.Input{
			Message: "Package Name:",
			Default: (func() string {
				pkgName, err := gelt.GetGoModuleName(".")
				if err != nil {
					return ""
				}

				return pkgName
			})(),
		},
		Validate: survey.Required,
	},
}

var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new gelt project.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		answers := newProjectArgs{}
		if len(args) > 0 {
			answers.ProjectDir = args[0]
		}

		err := survey.Ask(qs, &answers)
		if err != nil {
			color.Red(err.Error())
			return
		}

		err = newProject(answers)
		if err != nil {
			color.Red(err.Error())
			return
		}
	},
}

type newProjectArgs struct {
	ProjectDir  string `survey:"project_dir"`
	PackageName string `survey:"package_name"`
}

func generateFile(templatePath string, outputPath string, data any) error {
	tmpl, err := template.New(templatePath).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, data)
	if err != nil {
		return err
	}

	content := builder.String()
	err = os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return err
	}
	defer os.Remove(templatePath)

	return nil
}

func newProject(args newProjectArgs) error {
	err := cp.Copy("_template", args.ProjectDir)
	if err != nil {
		return err
	}

	err = os.Chdir(args.ProjectDir)
	if err != nil {
		return err
	}

	err = generateFile(
		"go.mod.tmpl",
		"go.mod",
		map[string]string{"ModuleName": args.PackageName},
	)
	if err != nil {
		return err
	}

	err = generateFile(
		"routes_gen.go.tmpl",
		"routes_gen.go",
		map[string]string{"ModuleName": args.PackageName},
	)
	if err != nil {
		return err
	}

	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		return err
	}

	color.Green("Project created successfully! 🎉")
	color.Blue(`
To run the project:

cd %s
go run main.go
	`, args.ProjectDir)

	return nil
}
