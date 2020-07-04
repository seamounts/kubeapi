package templates

import (
	"github.com/seamounts/kubeapi/pkg/model/file"
)

type GoMod struct {
	file.TemplateMixin
	file.RepositoryMixin
}

// GetBody implements Template
func (f *GoMod) GetBody() string {
	return f.TemplateBody
}

func (f *GoMod) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = "go.mod"
	}

	f.TemplateBody = goModTemplae
	f.IfExistsAction = file.Skip

	return nil
}

const goModTemplae = `
module {{ .Repo }}

go 1.13
`
