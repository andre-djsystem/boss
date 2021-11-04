package uninstall

import (
	"path/filepath"

	"github.com/hashload/boss/internal/pkg/configuration"
	"github.com/hashload/boss/internal/pkg/models"
	"github.com/hashload/boss/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdUnstall add the command line uninstall
func NewCmdUnstall(config *configuration.Configuration) *cobra.Command {
	var noSave bool
	cmd := &cobra.Command{
		Use:     "uninstall",
		Short:   "Uninstall a dependency",
		Long:    "This uninstalls a package, completely removing everything boss installed on its behalf",
		Aliases: []string{"remove", "rm", "r", "un", "unlink"},
		Example: `  Uninstall a package:
  boss uninstall <pkg>

  Uninstall a package without removing it from the boss.json file:
  boss uninstall <pkg> --no-save`,
		Run: func(cmd *cobra.Command, args []string) {
			err := uninstallDependency(config, noSave, args)
			util.CheckErr(err)
		},
	}
	cmd.Flags().BoolVar(&noSave, "no-save", false, "package will not be removed from your boss.json file")
	return cmd
}

func uninstallDependency(config *configuration.Configuration, noSave bool, args []string) error {
	currentDir, err := config.CurrentDir()
	if err != nil {
		return err
	}

	bossPath := filepath.Join(currentDir, "boss.json")
	pkg, err := models.LoadPackage(bossPath)
	if err != nil {
		return err
	}

	for dependency := range args {
		dependencyRepository := util.ParseDependency(args[dependency])
		//TODO implement remove without reinstall process
		pkg.UninstallDependency(dependencyRepository)
	}

	_, err = pkg.SaveToFile(bossPath)
	return err
}
