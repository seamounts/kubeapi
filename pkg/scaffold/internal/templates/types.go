package templates

import (
	"path/filepath"

	"github.com/seamounts/kubeapi/pkg/model/file"
)

type Types struct {
	file.TemplateMixin
	file.ResourceMixin
	file.RepositoryMixin
}

// GetBody implements Template
func (t *Types) GetBody() string {
	return t.TemplateBody
}

func (f *Types) SetTemplateDefaults() error {
	f.Path = filepath.Join("apis", "%[group]", "%[version]", "%[kind]_types.go")
	f.Path = f.Resource.Replacer().Replace(f.Path)

	f.TemplateBody = typesTemplate

	f.IfExistsAction = file.Overwrite

	return nil
}

const typesTemplate = `
package {{ .Resource.Version }}

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// {{ .Resource.Kind }} is a specification for a {{ .Resource.Kind }} resource
type {{ .Resource.Kind }} struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `

	Spec   {{ .Resource.Kind }}Spec    ` + "`" + `json:"spec,omitempty"` + "`" + `
	Status {{ .Resource.Kind }}Status ` + "`" + `json:"status,omitempty"` + "`" + `
}

// {{ .Resource.Kind }}Spec is the spec for a {{ .Resource.Kind }} resource
type {{ .Resource.Kind }}Spec struct {
	// Foo is an example field of {{ .Resource.Kind }}. Edit {{ .Resource.Kind }}_types.go to remove/update
	Foo string ` + "`" + `json:"foo,omitempty"` + "`" + `
}

// {{ .Resource.Kind }}Status is the status for a {{ .Resource.Kind }} resource
type {{ .Resource.Kind }}Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// {{ .Resource.Kind }}List is a list of {{ .Resource.Kind }} resources
type {{ .Resource.Kind }}List struct {
	metav1.TypeMeta ` + "`" + `json:",inline"` + "`" + `
	metav1.ListMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
	Items           []{{ .Resource.Kind }} ` + "`" + `json:"items"` + "`" + `
}

func init() {
	SchemeBuilder.Register(&{{ .Resource.Kind }}{}, &{{ .Resource.Kind }}List{})
}
`
