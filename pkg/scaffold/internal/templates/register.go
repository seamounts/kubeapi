package templates

import (
	"path/filepath"

	"github.com/seamounts/kubeapi/pkg/model/file"
)

type Register struct {
	file.TemplateMixin
	file.ResourceMixin
	file.RepositoryMixin
}

// GetBody implements Template
func (f *Register) GetBody() string {
	return f.TemplateBody
}

func (f *Register) SetTemplateDefaults() error {
	f.Path = filepath.Join("apis", "%[group]", "%[version]", "register.go")
	f.Path = f.Resource.Replacer().Replace(f.Path)

	f.TemplateBody = registerTemplate

	f.IfExistsAction = file.Overwrite

	return nil
}

const registerTemplate = `
package  {{ .Resource.Version }}

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: Version}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder initializes a scheme builder
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)
`
