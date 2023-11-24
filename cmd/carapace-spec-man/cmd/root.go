package cmd

import (
	"fmt"

	"github.com/rsteube/carapace"
	spec "github.com/rsteube/carapace-spec-man"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-spec-man <executable>",
	Short: "generate spec from manpages",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command, err := spec.Command(args[0], !cmd.Flag("no-trim").Changed)
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

	rootCmd.Flags().Bool("no-trim", false, "don't trim descriptions")

	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionExecutables(),
	)
}
