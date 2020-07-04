package v1

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/seamounts/kubeapi/internal/cmdutil"
	"github.com/seamounts/kubeapi/pkg/codegen"
	"github.com/seamounts/kubeapi/pkg/model/config"
	"github.com/seamounts/kubeapi/pkg/model/resource"
	"github.com/seamounts/kubeapi/pkg/plugin"
	"github.com/seamounts/kubeapi/pkg/scaffold"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

type createAPIPlugin struct {
	config *config.Config

	resource *resource.Options

	// force indicates that the resource should be created even if it already exists
	force bool
}

var (
	_ plugin.CreateAPI   = &createAPIPlugin{}
	_ cmdutil.RunOptions = &createAPIPlugin{}
)

func (p createAPIPlugin) UpdateContext(ctx *plugin.Context) {
	ctx.Description = `Scaffold a Kubernetes API by creating a Resource definition and / or a Controller.

create resource will prompt the user for if it should scaffold the Resource and / or Controller.  To only
scaffold a Controller for an existing Resource, select "n" for Resource.  To only define
the schema for a Resource without writing a Controller, select "n" for Controller.

After the scaffold is written, api will run make on the project.
`
	ctx.Examples = fmt.Sprintf(`  # Create a frigates API with Group: ship, Version: v1beta1 and Kind: Frigate
  %s create api --group ship --version v1beta1 --kind Frigate

  # Edit the API Scheme
  nano api/v1beta1/frigate_types.go

  # Edit the Controller
  nano controllers/frigate/frigate_controller.go

  # Edit the Controller Test
  nano controllers/frigate/frigate_controller_test.go

  # Install CRDs into the Kubernetes cluster using kubectl apply
  make install

  # Regenerate code and run against the Kubernetes cluster configured by ~/.kube/config
  make run
	`,
		ctx.CommandName)
}

func (p *createAPIPlugin) BindFlags(fs *pflag.FlagSet) {

	fs.BoolVar(&p.force, "force", false,
		"attempt to create resource even if it already exists")

	p.resource = &resource.Options{}
	fs.StringVar(&p.resource.Kind, "kind", "", "resource Kind")
	fs.StringVar(&p.resource.Group, "group", "", "resource Group")
	fs.StringVar(&p.resource.Version, "version", "", "resource Version")
	fs.BoolVar(&p.resource.Namespaced, "namespaced", true, "resource is namespaced")
}

func (p *createAPIPlugin) InjectConfig(c *config.Config) {
	p.config = c
}

func (p *createAPIPlugin) Run() error {
	return cmdutil.Run(p)
}

func (p *createAPIPlugin) Validate() error {
	if err := p.resource.Validate(); err != nil {
		return err
	}

	// Check that resource doesn't exist or flag force was set
	if !p.force && p.config.HasResource(p.resource.GVK()) {
		return errors.New("API resource already exists")
	}

	return nil
}

func (p *createAPIPlugin) GetScaffolder() (scaffold.Scaffolder, error) {
	// Load the boilerplate
	bp, err := ioutil.ReadFile(filepath.Join("hack", "boilerplate.go.txt")) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("unable to load boilerplate: %v", err)
	}

	// Create the actual resource from the resource options
	res := p.resource.NewResource(p.config)
	return scaffold.NewAPIScaffolder(p.config, string(bp), res), nil
}

func (p *createAPIPlugin) PostScaffold() error {
	klog.Infoln("Start Generating Client")
	err := codegen.GetCodeGen(p.config, p.resource).Run()
	if err != nil {
		return err
	}
	return nil
}
