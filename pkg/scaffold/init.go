package scaffold

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/seamounts/kubeapi/pkg/model"
	"github.com/seamounts/kubeapi/pkg/model/config"
	"github.com/seamounts/kubeapi/pkg/scaffold/internal/machinery"
	"github.com/seamounts/kubeapi/pkg/scaffold/internal/templates"
)

type initScaffolder struct {
	config          *config.Config
	boilerplatePath string
	license         string
	owner           string
}

// NewInitScaffolder returns a new Scaffolder for project initialization operations
func NewInitScaffolder(config *config.Config, license, owner string) Scaffolder {
	return &initScaffolder{
		config:          config,
		boilerplatePath: filepath.Join("hack", "boilerplate.go.txt"),
		license:         license,
		owner:           owner,
	}
}

func (s *initScaffolder) newUniverse(boilerplate string) *model.Universe {
	return model.NewUniverse(
		model.WithConfig(s.config),
		model.WithBoilerplate(boilerplate),
	)
}

// Scaffold implements Scaffolder
func (s *initScaffolder) Scaffold() error {
	fmt.Println("Writing scaffold for you to edit...")

	switch {
	case s.config.IsV1():
		return s.scaffold()
	default:
		return fmt.Errorf("unknown project version %v", s.config.Version)
	}
}

func (s *initScaffolder) scaffold() error {
	bpFile := &templates.Boilerplate{}
	bpFile.Path = s.boilerplatePath
	bpFile.License = s.license
	bpFile.Owner = s.owner

	if err := machinery.NewScaffold().Execute(
		s.newUniverse(""),
		bpFile,
	); err != nil {
		return err
	}

	boilerplate, err := ioutil.ReadFile(s.boilerplatePath) //nolint:gosec
	if err != nil {
		return err
	}

	return machinery.NewScaffold().Execute(
		s.newUniverse(string(boilerplate)),
		&templates.GitIgnore{},
		&templates.GoMod{},
	)

	return nil
}
