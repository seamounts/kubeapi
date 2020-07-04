package file

import (
	"text/template"

	"github.com/seamounts/kubeapi/pkg/model/resource"
)

// Builder defines the basic methods that any file builder must implement
type Builder interface {
	// GetPath returns the path to the file location
	GetPath() string
	// GetIfExistsAction returns the behavior when creating a file that already exists
	GetIfExistsAction() IfExistsAction
}

// Template is file builder based on a file template
type Template interface {
	Builder
	// GetBody returns the template body
	GetBody() string
	// SetTemplateDefaults sets the default values for templates
	SetTemplateDefaults() error
}

// HasRepository allows the repository to be used on a template
type HasRepository interface {
	// InjectRepository sets the template repository
	InjectRepository(string)
}

// HasResource allows a resource to be used on a template
type HasResource interface {
	// InjectResource sets the template resource
	InjectResource(*resource.Resource)
}

// UseCustomFuncMap allows a template to use a custom template.FuncMap instead of the default FuncMap.
type UseCustomFuncMap interface {
	// GetFuncMap returns a custom FuncMap.
	GetFuncMap() template.FuncMap
}
