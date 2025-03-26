package main

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "gelt",
	Short: "Gelt CLI",
}

func main() {
	RootCmd.AddCommand(NewCmd)
	RootCmd.Execute()
}
