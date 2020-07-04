package deepcopy

import (
	"path/filepath"

	generatorargs "k8s.io/code-generator/cmd/deepcopy-gen/args"
	"k8s.io/code-generator/pkg/util"
	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/deepcopy-gen/generators"
	"k8s.io/klog"
)

type OptionsFunc func(genericArgs *args.GeneratorArgs, customArgs *generatorargs.CustomArgs) error

type DeepCopy struct {
	genericArgs *args.GeneratorArgs
}

func NewDeepCopy(option OptionsFunc) (*DeepCopy, error) {
	genericArgs, customArgs := generatorargs.NewDefaults()

	// Override defaults.
	// TODO: move this out of deepcopy-gen
	genericArgs.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), util.BoilerplatePath())

	option(genericArgs, customArgs)
	if err := generatorargs.Validate(genericArgs); err != nil {
		return nil, err
	}

	dc := &DeepCopy{
		genericArgs: genericArgs,
	}

	return dc, nil
}

func (dc *DeepCopy) Run() error {
	// Run it.
	if err := dc.genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		return err
	}
	klog.V(2).Info("Completed successfully.")

	return nil
}
