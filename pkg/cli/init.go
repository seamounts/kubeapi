package cli

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	internalconfig "github.com/seamounts/kubeapi/internal/config"
	"github.com/seamounts/kubeapi/pkg/plugin"
	"github.com/spf13/cobra"
)

func (c *cli) newInitCmd() *cobra.Command {
	ctx := c.newInitContext()
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialize a new project",
		Long:    ctx.Description,
		Example: ctx.Examples,
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	// Register --project-version on the dynamically created command
	// so that it shows up in help and does not cause a parse error.
	cmd.Flags().String(projectVersionFlag, c.defaultProjectVersion,
		fmt.Sprintf("project version, possible values: (%s)", strings.Join(c.getAvailableProjectVersions(), ", ")))

	// If only the help flag was set, return the command as is.
	if c.doGenericHelp {
		return cmd
	}

	// Lookup the plugin for projectVersion and bind it to the command.
	c.bindInit(ctx, cmd)
	return cmd
}

func (c cli) newInitContext() plugin.Context {
	return plugin.Context{
		CommandName: c.commandName,
		Description: `Initialize a new project.

For further help about a specific project version, set --project-version.
`,
		Examples: c.getInitHelpExamples(),
	}
}

func (c cli) getInitHelpExamples() string {
	var sb strings.Builder
	for _, version := range c.getAvailableProjectVersions() {
		rendered := fmt.Sprintf(`  # Help for initializing a project with version %s
  %s init --project-version=%s -h

`,
			version, c.commandName, version)
		sb.WriteString(rendered)
	}
	return strings.TrimSuffix(sb.String(), "\n\n")
}

func (c cli) getAvailableProjectVersions() (projectVersions []string) {
	versionSet := make(map[string]struct{})
	for version, _ := range c.plugins {
		versionSet[version] = struct{}{}
	}
	for version := range versionSet {
		projectVersions = append(projectVersions, strconv.Quote(version))
	}
	return projectVersions
}

func (c cli) bindInit(ctx plugin.Context, cmd *cobra.Command) {
	getter, isGetter := c.resolvedPlugin.(plugin.InitPluginGetter)
	if getter == nil || !isGetter {
		err := fmt.Errorf("plugin does not support an Init plugin")
		cmdErr(cmd, err)
		return
	}

	cfg := internalconfig.New(internalconfig.DefaultPath)
	cfg.Version = c.projectVersion

	init := getter.GetInitPlugin()
	init.InjectConfig(&cfg.Config)
	init.BindFlags(cmd.Flags())
	init.UpdateContext(&ctx)
	cmd.Long = ctx.Description
	cmd.Example = ctx.Examples
	cmd.RunE = func(*cobra.Command, []string) error {
		// Check if a config is initialized in the command runner so the check
		// doesn't erroneously fail other commands used in initialized projects.
		_, err := internalconfig.Read()
		if err == nil || os.IsExist(err) {
			log.Fatal("config already initialized")
		}
		if err := init.Run(); err != nil {
			return fmt.Errorf("failed to initialize project with version %q: %v", c.projectVersion, err)
		}
		return cfg.Save()
	}
}
