package cli

import (
	"fmt"
	"log"
	"os"

	internalconfig "github.com/seamounts/kubeapi/internal/config"
	"github.com/seamounts/kubeapi/pkg/internal/validation"
	"github.com/seamounts/kubeapi/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	runInProjectRootMsg = `For project-specific information, run this command in the root directory of a
project.
`

	projectVersionFlag = "project-version"
	helpFlag           = "help"
)

// CLI interacts with a command line interface.
type CLI interface {
	// Run runs the CLI, usually returning an error if command line configuration
	// is incorrect.
	Run() error
}

// Option is a function that can configure the cli
type Option func(*cli) error

// cli defines the command line structure and interfaces that are used to
// scaffold kubebuilder project files.
type cli struct {
	// Base command name. Can be injected downstream.
	commandName string
	// Default project version. Used in CLI flag setup.
	defaultProjectVersion string
	// Project version to scaffold.
	projectVersion string
	// True if the project has config file.
	configured bool
	// Whether the command is requesting help.
	doGenericHelp bool

	defaultPlugin plugin.Base

	// Plugins injected by options.
	plugins map[string]plugin.Base

	resolvedPlugin plugin.Base

	// Base command.
	cmd *cobra.Command
	// Commands injected by options.
	extraCommands []*cobra.Command
}

func New(opts ...Option) (CLI, error) {
	cli := &cli{
		commandName:           "kubeapi",
		defaultProjectVersion: internalconfig.DefaultVersion,
		plugins:               make(map[string]plugin.Base),
	}

	for _, opt := range opts {
		if err := opt(cli); err != nil {
			return nil, err
		}
	}

	cli.initialize()

	return cli, nil
}

// Run runs the cli.
func (c cli) Run() error {
	return c.cmd.Execute()
}

// WithCommandName is an Option that sets the cli's root command name.
func WithCommandName(name string) Option {
	return func(c *cli) error {
		c.commandName = name
		return nil
	}
}

// WithDefaultProjectVersion is an Option that sets the cli's default project
// version. Setting an unknown version will result in an error.
func WithDefaultProjectVersion(version string) Option {
	return func(c *cli) error {
		if err := validation.ValidateProjectVersion(version); err != nil {
			return fmt.Errorf("broken pre-set default project version %q: %v", version, err)
		}
		c.defaultProjectVersion = version
		return nil
	}
}

// WithPlugins is an Option that sets the cli's plugins.
func WithPlugins(plugins ...plugin.Base) Option {
	return func(c *cli) error {
		for _, p := range plugins {
			for _, version := range p.SupportedProjectVersions() {
				c.plugins[version] = p
			}
		}
		for _, p := range c.plugins {
			if err := validatePlugins(p); err != nil {
				return fmt.Errorf("broken pre-set plugins: %v", err)
			}
		}
		return nil
	}
}

// WithDefaultPlugins is an Option that sets the cli's default plugins. Only
// one plugin per project version is allowed.
func WithDefaultPlugin(p plugin.Base) Option {
	return func(c *cli) error {
		if err := validatePlugin(p); err != nil {
			return fmt.Errorf("broken pre-set default plugin %q: %v", plugin.KeyFor(p), err)
		}
		c.defaultPlugin = p

		return nil
	}
}

// WithExtraCommands is an Option that adds extra subcommands to the cli.
// Adding extra commands that duplicate existing commands results in an error.
func WithExtraCommands(cmds ...*cobra.Command) Option {
	return func(c *cli) error {
		c.extraCommands = append(c.extraCommands, cmds...)
		return nil
	}
}

func (c *cli) initialize() error {

	// Initialize cli with globally-relevant flags or flags that determine
	// certain plugin type's configuration.
	if err := c.parseBaseFlags(); err != nil {
		return err
	}

	// Configure the project version first for plugin retrieval in command
	// constructors.
	projectConfig, err := internalconfig.Read()
	if os.IsNotExist(err) {
		c.configured = false
		if c.projectVersion == "" {
			c.projectVersion = c.defaultProjectVersion
		}
	} else if err == nil {
		c.configured = true
		c.projectVersion = projectConfig.Version
	} else {
		return fmt.Errorf("failed to read config: %v", err)
	}

	// Validate after setting projectVersion but before buildRootCmd so we error
	// out before an error resulting from an incorrect cli is returned downstream.
	if err = c.validate(); err != nil {
		return err
	}

	c.resolvedPlugin = c.plugins[c.projectVersion]
	if c.resolvedPlugin == nil {
		c.resolvedPlugin = c.defaultPlugin
	}

	c.cmd = c.buildRootCmd()
	// Add extra commands injected by options.
	for _, cmd := range c.extraCommands {
		for _, subCmd := range c.cmd.Commands() {
			if cmd.Name() == subCmd.Name() {
				return fmt.Errorf("command %q already exists", cmd.Name())
			}
		}
		c.cmd.AddCommand(cmd)
	}

	return nil
}

// parseBaseFlags parses the command line arguments, looking for flags that
// affect initialization of a cli. An error is returned only if an error
// unrelated to flag parsing occurs.
func (c *cli) parseBaseFlags() error {
	// Create a dummy "base" flagset to populate from CLI args.
	fs := pflag.NewFlagSet("base", pflag.ExitOnError)
	fs.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}

	var help bool
	// Set base flags that require pre-parsing to initialize c.
	fs.BoolVarP(&help, helpFlag, "h", false, "print help")
	fs.StringVar(&c.projectVersion, projectVersionFlag, c.defaultProjectVersion, "project version")

	// Parse current CLI args outside of cobra.
	err := fs.Parse(os.Args[1:])
	// User needs *generic* help if args are incorrect or --help is set and
	// --project-version is not set. Plugin-specific help is given if a
	// plugin.Context is updated, which does not require this field.
	c.doGenericHelp = err != nil || help && !fs.Lookup(projectVersionFlag).Changed

	return nil
}

// validate validates fields in a cli.
func (c cli) validate() error {
	// Validate project version.
	if err := validation.ValidateProjectVersion(c.projectVersion); err != nil {
		return fmt.Errorf("invalid project version %q: %v", c.projectVersion, err)
	}

	if _, versionFound := c.plugins[c.projectVersion]; !versionFound {
		return fmt.Errorf("no plugins for project version %q", c.projectVersion)
	}

	// Validate plugin versions and name.
	for _, versionedPlugin := range c.plugins {
		if err := validatePlugins(versionedPlugin); err != nil {
			return err
		}
	}

	return nil
}

// validatePlugins validates the name and versions of a list of plugins.
func validatePlugins(plugins ...plugin.Base) error {
	pluginNameSet := make(map[string]struct{}, len(plugins))
	for _, p := range plugins {
		if err := validatePlugin(p); err != nil {
			return err
		}
		// Check for duplicate plugin keys.
		pluginKey := plugin.KeyFor(p)
		if _, seen := pluginNameSet[pluginKey]; seen {
			return fmt.Errorf("two plugins have the same key: %q", pluginKey)
		}
		pluginNameSet[pluginKey] = struct{}{}
	}
	return nil
}

// validatePlugin validates the name and versions of a plugin.
func validatePlugin(p plugin.Base) error {
	pluginName := p.Name()
	if err := plugin.ValidateName(pluginName); err != nil {
		return fmt.Errorf("invalid plugin name %q: %v", pluginName, err)
	}
	pluginVersion := p.Version()
	if err := plugin.ValidateVersion(pluginVersion); err != nil {
		return fmt.Errorf("invalid plugin %q version %q: %v",
			pluginName, pluginVersion, err)
	}
	for _, projectVersion := range p.SupportedProjectVersions() {
		if err := validation.ValidateProjectVersion(projectVersion); err != nil {
			return fmt.Errorf("invalid plugin %q supported project version %q: %v",
				pluginName, projectVersion, err)
		}
	}
	return nil
}

// buildRootCmd returns a root command with a subcommand tree reflecting the
// current project's state.
func (c cli) buildRootCmd() *cobra.Command {
	rootCmd := c.defaultCommand()

	// kubebuilder create
	createCmd := c.newCreateCmd()
	// kubebuilder create api
	createCmd.AddCommand(c.newCreateAPICmd())
	if createCmd.HasSubCommands() {
		rootCmd.AddCommand(createCmd)
	}

	// kubebuilder init
	rootCmd.AddCommand(c.newInitCmd())

	return rootCmd
}

// defaultCommand returns the root command without its subcommands.
func (c cli) defaultCommand() *cobra.Command {
	return &cobra.Command{
		Use:   c.commandName,
		Short: "Development kit for building Kubernetes extensions and tools.",
		Long: fmt.Sprintf(`Development kit for building Kubernetes extensions and tools.

Provides libraries and tools to create new projects, APIs and controllers.
Includes tools for packaging artifacts into an installer container.

Typical project lifecycle:

- initialize a project:

  %s init --license apache2 --owner "The Kubernetes authors"

- create one or more a new resource APIs and add your code to them:

  %s create api --group <group> --version <version> --kind <Kind>
After the scaffold is written, api will run make on the project.
`,
			c.commandName, c.commandName),
		Example: fmt.Sprintf(`
  # Initialize your project
  %s init --license apache2 --owner "The Kubernetes authors"

  # Create a frigates API with Group: ship, Version: v1beta1 and Kind: Frigate
  %s create api --group ship --version v1beta1 --kind Frigate

  # Edit the API Scheme
  nano api/v1beta1/frigate_types.go
`,
			c.commandName, c.commandName),

		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
	}
}
