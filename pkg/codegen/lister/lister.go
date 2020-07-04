package lister

import (
	"path/filepath"

	"k8s.io/code-generator/cmd/lister-gen/generators"
	"k8s.io/code-generator/pkg/util"
	"k8s.io/gengo/args"
	"k8s.io/klog"

	generatorargs "k8s.io/code-generator/cmd/lister-gen/args"
)

type OptionsFunc func(genericArgs *args.GeneratorArgs, customArgs *generatorargs.CustomArgs) error

type Lister struct {
	genericArgs *args.GeneratorArgs
}

func NewLister(opt OptionsFunc) (*Lister, error) {
	genericArgs, customArgs := generatorargs.NewDefaults()

	// Override defaults.
	// TODO: move this out of lister-gen
	genericArgs.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), util.BoilerplatePath())
	genericArgs.OutputPackagePath = "k8s.io/kubernetes/pkg/client/listers"

	opt(genericArgs, customArgs)
	if err := generatorargs.Validate(genericArgs); err != nil {
		return nil, err
	}

	return &Lister{
		genericArgs: genericArgs,
	}, nil
}

func (l *Lister) Run() error {
	// Run it.
	if err := l.genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		return err
	}
	klog.V(2).Info("Completed successfully.")

	return nil
}
