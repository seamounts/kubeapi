package v1

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/seamounts/kubeapi/internal/cmdutil"
	"github.com/seamounts/kubeapi/pkg/internal/validation"
	"github.com/seamounts/kubeapi/pkg/model/config"
	"github.com/seamounts/kubeapi/pkg/plugin"
	"github.com/spf13/pflag"

	"github.com/seamounts/kubeapi/pkg/plugin/internal"
	"github.com/seamounts/kubeapi/pkg/scaffold"
)

type initPlugin struct {
	config *config.Config
	// For help text.
	commandName string

	// boilerplate options
	license string
	owner   string

	// flags
	skipGoVersionCheck bool
}

var (
	_ plugin.Init        = &initPlugin{}
	_ cmdutil.RunOptions = &initPlugin{}
)

func (p *initPlugin) UpdateContext(ctx *plugin.Context) {
	ctx.Description = `Initialize a new project including vendor/ directory and Go package directories.

Writes the following files:
- a boilerplate license file
- a PROJECT file with the domain and repo
- a go.mod with project dependencies

`
	ctx.Examples = fmt.Sprintf(`  # Scaffold a project using the apache2 license with "The Kubernetes authors" as owners
  %s init --project-version=1 --domain example.org --license apache2 --owner "The Kubernetes authors"
`,
		ctx.CommandName)

	p.commandName = ctx.CommandName
}

func (p *initPlugin) BindFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&p.skipGoVersionCheck, "skip-go-version-check",
		false, "if specified, skip checking the Go version")

	// boilerplate args
	fs.StringVar(&p.license, "license", "apache2",
		"license to use to boilerplate, may be one of 'apache2', 'none'")
	fs.StringVar(&p.owner, "owner", "", "owner to add to the copyright")

	// project args
	fs.StringVar(&p.config.Repo, "repo", "", "name to use for go module (e.g., github.com/user/repo), "+
		"defaults to the go package of the current working directory.")
	fs.StringVar(&p.config.Domain, "domain", "my.domain", "domain for groups")
}

func (p *initPlugin) InjectConfig(c *config.Config) {
	p.config = c
}

func (p *initPlugin) Run() error {
	return cmdutil.Run(p)
}

func (p *initPlugin) Validate() error {
	// Requires go1.11+
	if !p.skipGoVersionCheck {
		if err := internal.ValidateGoVersion(); err != nil {
			return err
		}
	}

	// Check if the project name is a valid namespace according to k8s
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error to get the current path: %v", err)
	}
	projectName := filepath.Base(dir)
	if err := validation.IsDNS1123Label(strings.ToLower(projectName)); err != nil {
		return fmt.Errorf("project name (%s) is invalid: %v", projectName, err)
	}

	// Try to guess repository if flag is not set.
	if p.config.Repo == "" {
		repoPath, err := internal.FindCurrentRepo()
		if err != nil {
			return fmt.Errorf("error finding current repository: %v", err)
		}
		p.config.Repo = repoPath
	}

	return nil
}

func (p *initPlugin) GetScaffolder() (scaffold.Scaffolder, error) {
	return scaffold.NewInitScaffolder(p.config, p.license, p.owner), nil
}

func (p *initPlugin) PostScaffold() error {
	fmt.Printf("Next: define a resource with:\n$ %s create api\n", p.commandName)
	return nil
}
