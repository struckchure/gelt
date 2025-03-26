package main

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
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
		Name:     "project_name",
		Prompt:   &survey.Input{Message: "Project Name:"},
		Validate: survey.Required,
		Transform: survey.TransformString(
			func(s string) string {
				return strings.ReplaceAll(s, " ", "_")
			},
		),
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
		answers := struct {
			ProjectDir  string `survey:"project_dir"`
			ProjectName string `survey:"project_name"`
			PackageName string `survey:"package_name"`
		}{}
		if len(args) > 0 {
			answers.ProjectDir = args[0]
		}

		err := survey.Ask(qs, &answers)
		if err != nil {
			color.Red(err.Error())
			return
		}
		color.Green("%#v\n", answers)
	},
}
