package informar

import (
	"path/filepath"

	"k8s.io/code-generator/cmd/informer-gen/generators"
	"k8s.io/code-generator/pkg/util"
	"k8s.io/gengo/args"
	"k8s.io/klog"

	generatorargs "k8s.io/code-generator/cmd/informer-gen/args"
)

type OptionsFunc func(genericArgs *args.GeneratorArgs, customArgs *generatorargs.CustomArgs) error

type Informar struct {
	genericArgs *args.GeneratorArgs
}

func NewInformar(opt OptionsFunc) (*Informar, error) {
	genericArgs, customArgs := generatorargs.NewDefaults()
	// Override defaults.
	// TODO: move out of informer-gen
	genericArgs.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), util.BoilerplatePath())
	genericArgs.OutputPackagePath = "k8s.io/kubernetes/pkg/client/informers/informers_generated"
	customArgs.VersionedClientSetPackage = "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
	customArgs.InternalClientSetPackage = "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	customArgs.ListersPackage = "k8s.io/kubernetes/pkg/client/listers"

	opt(genericArgs, customArgs)
	if err := generatorargs.Validate(genericArgs); err != nil {
		return nil, err
	}

	return &Informar{
		genericArgs: genericArgs,
	}, nil
}

func (in *Informar) Run() error {
	// Run it.
	if err := in.genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		return nil
	}
	klog.V(2).Info("Completed successfully.")

	return nil
}
