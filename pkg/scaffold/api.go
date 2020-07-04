package scaffold

import (
	"fmt"

	"github.com/seamounts/kubeapi/pkg/model"
	"github.com/seamounts/kubeapi/pkg/model/config"
	"github.com/seamounts/kubeapi/pkg/model/resource"
	"github.com/seamounts/kubeapi/pkg/scaffold/internal/machinery"
	"github.com/seamounts/kubeapi/pkg/scaffold/internal/templates"
)

type apiScaffolder struct {
	config      *config.Config
	resource    *resource.Resource
	boilerplate string
}

func NewAPIScaffolder(config *config.Config, boilerplate string, res *resource.Resource) Scaffolder {
	s := &apiScaffolder{
		config:   config,
		resource: res,
	}

	return s
}

// Scaffold implements Scaffolder
func (s *apiScaffolder) Scaffold() error {
	fmt.Println("Writing scaffold for you to edit...")
	return s.scaffold()
}

func (s *apiScaffolder) scaffold() error {
	machinery.NewScaffold().Execute(
		s.newUniverse(),
		&templates.Types{},
		&templates.Doc{},
		&templates.Register{},
	)

	return nil
}

func (s *apiScaffolder) newUniverse() *model.Universe {
	return model.NewUniverse(
		model.WithConfig(s.config),
		model.WithResource(s.resource),
	)

}
