package plugin

import (
	"path"
	"strings"

	"github.com/seamounts/kubeapi/pkg/model/config"
	"github.com/spf13/pflag"
)

const DefaultNameQualifier = ".kubeapi.io"

type Base interface {
	// Name returns a DNS1123 label string defining the plugin type.
	// For example, kubeapi's main plugin would return "go".
	Name() string
	// Version returns the plugin's semantic version, ex. "v1.2.3".
	//
	// Note: this version is different from config version.
	Version() string
	// SupportedProjectVersions lists all project configuration versions this
	// plugin supports, ex. []string{"2", "3"}. The returned slice cannot be empty.
	SupportedProjectVersions() []string
}

// Key returns a unique identifying string for a plugin's name and version.
func Key(name, version string) string {
	if version == "" {
		return name
	}
	return path.Join(name, "v"+strings.TrimLeft(version, "v"))
}

// KeyFor returns a Base plugin's unique identifying string.
func KeyFor(p Base) string {
	return Key(p.Name(), p.Version())
}

// SplitKey returns a name and version for a plugin key.
func SplitKey(key string) (string, string) {
	if !strings.Contains(key, "/") {
		return key, ""
	}
	keyParts := strings.SplitN(key, "/", 2)
	return keyParts[0], keyParts[1]
}

// GetShortName returns plugin's short name (name before domain) if name
// is fully qualified (has a domain suffix), otherwise GetShortName returns name.
func GetShortName(name string) string {
	return strings.SplitN(name, ".", 2)[0]
}

type GenericSubcommand interface {
	// UpdateContext updates a Context with command-specific help text, like description and examples.
	// Can be a no-op if default help text is desired.
	UpdateContext(*Context)
	// BindFlags binds the plugin's flags to the CLI. This allows each plugin to define its own
	// command line flags for the kubebuilder subcommand.
	BindFlags(*pflag.FlagSet)
	// Run runs the subcommand.
	Run() error
	// InjectConfig passes a config to a plugin. The plugin may modify the
	// config. Initializing, loading, and saving the config is managed by the
	// cli package.
	InjectConfig(*config.Config)
}

type Context struct {
	// CommandName sets the command name for a plugin.
	CommandName string
	// Description is a description of what this subcommand does. It is used to display help.
	Description string
	// Examples are one or more examples of the command-line usage
	// of this plugin's project subcommand support. It is used to display help.
	Examples string
}

type InitPluginGetter interface {
	Base
	// GetInitPlugin returns the underlying Init interface.
	GetInitPlugin() Init
}

type Init interface {
	GenericSubcommand
}

type CreateAPIPluginGetter interface {
	Base
	// GetCreateAPIPlugin returns the underlying CreateAPI interface.
	GetCreateAPIPlugin() CreateAPI
}

type CreateAPI interface {
	GenericSubcommand
}
