package codegen

import (
	"fmt"

	"github.com/seamounts/kubeapi/pkg/codegen/clientset"
	"github.com/seamounts/kubeapi/pkg/codegen/deepcopy"
	"github.com/seamounts/kubeapi/pkg/codegen/informar"
	"github.com/seamounts/kubeapi/pkg/codegen/lister"
	"github.com/seamounts/kubeapi/pkg/model/config"
	"github.com/seamounts/kubeapi/pkg/model/resource"
	clientsetargs "k8s.io/code-generator/cmd/client-gen/args"
	deepcopyaargs "k8s.io/code-generator/cmd/deepcopy-gen/args"
	informarargs "k8s.io/code-generator/cmd/informer-gen/args"
	listerargs "k8s.io/code-generator/cmd/lister-gen/args"
	"k8s.io/gengo/args"
	deepcopygenerators "k8s.io/gengo/examples/deepcopy-gen/generators"
	"k8s.io/klog/v2"
)

const (
	CLIENTSET_NAME_VERSIONED = "versioned"
	CLIENTSET_PKG_NAME       = "clientset"
	OUTPUT_DIR               = "client"
	INPUT_DIR                = "apis"
)

type CodeGen struct {
	config   *config.Config
	resource *resource.Resource

	inputDir  string
	outputDir string
}

var defaultCodeGen *CodeGen

func GetCodeGen(config *config.Config, opt *resource.Options) *CodeGen {
	if defaultCodeGen == nil {
		defaultCodeGen = &CodeGen{
			config:   config,
			resource: opt.NewResource(config),
		}
	}

	return defaultCodeGen
}

func (gen *CodeGen) Run() error {

	outputpkg := fmt.Sprintf("%s/%s", defaultCodeGen.config.Repo, OUTPUT_DIR)

	klog.Infoln("Generating deepcopy funcs")
	dc, err := deepcopy.NewDeepCopy(deepCopyOptions)
	if err != nil {
		return err
	}
	if err := dc.Run(); err != nil {
		return err
	}

	klog.Infof("Generating clientset for %s/%s at %s/%s", gen.resource.Group, gen.resource.Version,
		outputpkg, CLIENTSET_PKG_NAME)
	cs, err := clientset.NewClientSet(clientsetOptions)
	if err != nil {
		return err
	}
	if err := cs.Run(); err != nil {
		return err
	}

	klog.Infof("Generating listers for %s/%s at %s/listers", gen.resource.Group, gen.resource.Version, outputpkg)
	li, err := lister.NewLister(listerOptions)
	if err != nil {
		return err
	}
	if err := li.Run(); err != nil {
		return err
	}

	klog.Infof("Generating informers for %s/%s at %s/informers", gen.resource.Group, gen.resource.Version, outputpkg)
	in, err := informar.NewInformar(informarOptions)
	if err != nil {
		return err
	}
	if err := in.Run(); err != nil {
		return err
	}

	return nil
}

func deepCopyOptions(genericArgs *args.GeneratorArgs, customArgs *deepcopyaargs.CustomArgs) error {
	genericArgs.InputDirs = append(genericArgs.InputDirs, fmt.Sprintf("%s/%s/%s/%s",
		defaultCodeGen.config.Repo, INPUT_DIR, defaultCodeGen.resource.Group, defaultCodeGen.resource.Version))

	genericArgs.OutputFileBaseName = "zz_generated.deepcopy"
	genericArgs.CustomArgs = &deepcopygenerators.CustomArgs{
		BoundingDirs: []string{
			fmt.Sprintf("%s/%s", defaultCodeGen.config.Repo, INPUT_DIR),
		},
	}

	return nil
}

func clientsetOptions(genericArgs *args.GeneratorArgs, customArgs *clientsetargs.CustomArgs) error {
	customArgs.ClientsetName = CLIENTSET_NAME_VERSIONED
	genericArgs.OutputPackagePath = fmt.Sprintf("%s/%s/%s", defaultCodeGen.config.Repo,
		OUTPUT_DIR, CLIENTSET_PKG_NAME)

	gvPackages := clientsetargs.NewGVPackagesValue(clientsetargs.NewGroupVersionsBuilder(&customArgs.Groups), nil)

	gvPackages.Set(fmt.Sprintf("%s/%s/%s/%s", defaultCodeGen.config.Repo,
		INPUT_DIR, defaultCodeGen.resource.Group, defaultCodeGen.resource.Version))

	// add group version package as input dirs for gengo
	for _, pkg := range customArgs.Groups {
		for _, v := range pkg.Versions {
			genericArgs.InputDirs = append(genericArgs.InputDirs, v.Package)
		}
	}

	genericArgs.CustomArgs = customArgs

	return nil
}

func informarOptions(genericArgs *args.GeneratorArgs, customArgs *informarargs.CustomArgs) error {
	genericArgs.InputDirs = append(genericArgs.InputDirs, fmt.Sprintf("%s/%s/%s/%s",
		defaultCodeGen.config.Repo, INPUT_DIR, defaultCodeGen.resource.Group, defaultCodeGen.resource.Version))
	genericArgs.OutputPackagePath = fmt.Sprintf("%s/%s/informers", defaultCodeGen.config.Repo, OUTPUT_DIR)

	customArgs.VersionedClientSetPackage = fmt.Sprintf("%s/%s/%s/%s", defaultCodeGen.config.Repo,
		OUTPUT_DIR, CLIENTSET_PKG_NAME, CLIENTSET_NAME_VERSIONED)

	customArgs.ListersPackage = fmt.Sprintf("%s/%s/listers", defaultCodeGen.config.Repo, OUTPUT_DIR)

	genericArgs.CustomArgs = customArgs

	return nil
}

func listerOptions(genericArgs *args.GeneratorArgs, customArgs *listerargs.CustomArgs) error {
	genericArgs.InputDirs = append(genericArgs.InputDirs, fmt.Sprintf("%s/%s/%s/%s",
		defaultCodeGen.config.Repo, INPUT_DIR, defaultCodeGen.resource.Group, defaultCodeGen.resource.Version))
	genericArgs.OutputPackagePath = fmt.Sprintf("%s/%s/listers", defaultCodeGen.config.Repo, OUTPUT_DIR)

	return nil
}
