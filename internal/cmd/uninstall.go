package cmd

import (
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall <repository>",
	Short: "uninstall a previously installed application",
	RunE: func(cmd *cobra.Command, args []string) error {
		// validate the args
		if len(args) != 1 {
			return cmd.Usage()
		}

		// get the *package.Package for the provided URL
		pkg, err := packageForURL(args[0])
		if err != nil {
			return err
		}

		// get the installation directory for the package
		dir, err := installDir()
		if err != nil {
			return err
		}

		// uninstall the package if present
		return pkg.Uninstall(dir)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
