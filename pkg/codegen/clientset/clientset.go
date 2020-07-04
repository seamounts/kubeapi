package clientset

import (
	"path/filepath"

	generatorargs "k8s.io/code-generator/cmd/client-gen/args"
	"k8s.io/code-generator/cmd/client-gen/generators"
	"k8s.io/code-generator/pkg/util"
	"k8s.io/gengo/args"
	"k8s.io/klog"
)

type OptionsFunc func(genericArgs *args.GeneratorArgs, customArgs *generatorargs.CustomArgs) error

type ClientSet struct {
	genericArgs *args.GeneratorArgs
}

func NewClientSet(option OptionsFunc) (*ClientSet, error) {
	genericArgs, customArgs := generatorargs.NewDefaults()
	// Override defaults.
	// TODO: move this out of client-gen
	genericArgs.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), util.BoilerplatePath())
	genericArgs.OutputPackagePath = "k8s.io/kubernetes/pkg/client/clientset_generated/"

	option(genericArgs, customArgs)
	if err := generatorargs.Validate(genericArgs); err != nil {
		return nil, err
	}

	dc := &ClientSet{
		genericArgs: genericArgs,
	}

	return dc, nil
}

func (sc *ClientSet) Run() error {
	// Run it.
	if err := sc.genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		return nil
	}

	klog.V(2).Info("Completed successfully.")
	return nil
}
