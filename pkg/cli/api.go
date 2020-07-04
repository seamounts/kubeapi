package cli

import (
	"fmt"

	"github.com/seamounts/kubeapi/internal/config"
	"github.com/seamounts/kubeapi/pkg/plugin"
	"github.com/spf13/cobra"
)

func (c *cli) newCreateAPICmd() *cobra.Command {
	ctx := c.newAPIContext()
	cmd := &cobra.Command{
		Use:     "api",
		Short:   "Scaffold a Kubernetes API",
		Long:    ctx.Description,
		Example: ctx.Examples,
		RunE: errCmdFunc(
			fmt.Errorf("api subcommand requires an existing project"),
		),
	}

	// Lookup the plugin for projectVersion and bind it to the command.
	c.bindCreateAPI(ctx, cmd)
	return cmd
}

func (c cli) newAPIContext() plugin.Context {
	ctx := plugin.Context{
		CommandName: c.commandName,
		Description: `Scaffold a Kubernetes API.
`,
	}
	if !c.configured {
		ctx.Description = fmt.Sprintf("%s\n%s", ctx.Description, runInProjectRootMsg)
	}
	return ctx
}

func (c cli) bindCreateAPI(ctx plugin.Context, cmd *cobra.Command) {
	getter, isGetter := c.resolvedPlugin.(plugin.CreateAPIPluginGetter)
	if getter == nil || !isGetter {
		err := fmt.Errorf("plugin does not support an API creation plugin")
		cmdErr(cmd, err)
		return
	}

	cfg, err := config.LoadInitialized()
	if err != nil {
		cmdErr(cmd, err)
		return
	}

	createAPI := getter.GetCreateAPIPlugin()
	createAPI.InjectConfig(&cfg.Config)
	createAPI.BindFlags(cmd.Flags())
	createAPI.UpdateContext(&ctx)
	cmd.Long = ctx.Description
	cmd.Example = ctx.Examples
	cmd.RunE = runECmdFunc(cfg, createAPI,
		fmt.Sprintf("failed to create API with version %q", c.projectVersion))
}
