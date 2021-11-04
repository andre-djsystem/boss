package install

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashload/boss/internal/pkg/configuration"
	"github.com/hashload/boss/internal/pkg/models"
	"github.com/hashload/boss/pkg/consts"
	"github.com/hashload/boss/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdInstall add the command line install
func NewCmdInstall(config *configuration.Configuration) *cobra.Command {
	var noSave bool
	cmd := &cobra.Command{
		Use:     "install",
		Short:   "Install a new dependency",
		Aliases: []string{"i", "add"},
		Example: `  Add a new dependency:
  boss install <pkg>

  Add a new version-specific dependency:
  boss install <pkg>@<version>

  Install a dependency without add it from the boss.json file:
  boss install <pkg> --no-save`,
		Run: func(cmd *cobra.Command, args []string) {
			err := installDependency(config, noSave, args)
			util.CheckErr(err)
		},
	}
	cmd.Flags().BoolVar(&noSave, "no-save", false, "prevents saving to `dependencies`")
	return cmd
}

func installDependency(config *configuration.Configuration, noSave bool, args []string) error {
	currentDir, err := config.CurrentDir()
	if err != nil {
		return err
	}
	bossPath := filepath.Join(currentDir, "boss.json")
	pkg, err := models.LoadPackage(bossPath)
	if err != nil {
		return err
	}
	if config.Global {
		err = installGlobalDependency(pkg, args)
	} else {
		err = installLocalDependency(pkg, args)
	}
	if noSave || err != nil {
		return err
	}
	_, err = pkg.SaveToFile(bossPath)
	return err
}

func installLocalDependency(pkg *models.BossPackage, args []string) error {
	ensureDependencyOfArgs(pkg, args)
	doInstall(pkg)
	return nil
}

func installGlobalDependency(pkg *models.BossPackage, args []string) error {
	ensureDependencyOfArgs(pkg, args)
	doInstall(pkg)
	return nil
}

func ensureDependencyOfArgs(pkg *models.BossPackage, args []string) {
	for e := range args {
		dependency := util.ParseDependency(args[e])
		dependency = strings.ToLower(dependency)

		re := regexp.MustCompile(`(?m)(?:(?P<host>.*)(?::(?P<version>[\^~]?(?:(?:(?:[0-9]+)(?:\.[0-9]+)(?:\.[0-9]+)?))))$|(?P<host_only>.*))`)
		match := make(map[string]string)
		split := re.FindStringSubmatch(dependency)

		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				match[name] = split[i]
			}
		}
		var ver, dep string
		if len(match["version"]) == 0 {
			ver = consts.MinimalDependencyVersion
			dep = match["host_only"]
		} else {
			ver = match["version"]
			dep = match["host"]
		}

		if strings.HasSuffix(strings.ToLower(dep), ".git") {
			dep = dep[:len(dep)-4]
		}

		pkg.AddDependency(dep, ver)
	}
}

func doInstall(pkg *models.BossPackage) {
	fmt.Println("Installing modules...")
}
