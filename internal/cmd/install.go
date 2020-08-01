package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/git-pkg/gpkg/pkg/constants"
)

var installCmd = &cobra.Command{
	Use:   "install <repository>",
	Short: "install a new application",
	RunE: func(cmd *cobra.Command, args []string) error {
		// validate the arguments
		if len(args) != 1 {
			return cmd.Usage()
		}

		// get a *package.Package for the provided URL
		pkg, err := packageForURL(args[0])
		if err != nil {
			return err
		}

		// get the installation directory
		dir, err := installDir()
		if err != nil {
			return err
		}

		// install the package to the local system
		fmt.Printf("installing package %s to %s\n", pkg.URL, dir)
		if err := pkg.Install(dir, &constants.LatestVersion); err != nil {
			return err
		}
		fmt.Println("done!")

		// verify whether or not the installation bin directory is in the current $PATH
		if path := os.Getenv("PATH"); path != "" {
			found := false
			for _, d := range strings.Split(path, ":") {
				if d == fmt.Sprintf("%s/bin", dir) {
					found = true
				}
			}
			if !found {
				fmt.Printf("could not find %s/bin in your $PATH, make sure you export it in order to use installed packages!\n", dir)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
