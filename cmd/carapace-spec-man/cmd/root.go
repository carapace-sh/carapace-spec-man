package cmd

import (
	"fmt"

	"github.com/rsteube/carapace"
	spec "github.com/rsteube/carapace-spec-man"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-spec-man",
	Short: "generate spec from manpages",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command, err := spec.Command(args[0])
		if err != nil {
			return err
		}
		out, err := yaml.Marshal(command)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).Standalone()

	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionExecutables(),
	)
}
